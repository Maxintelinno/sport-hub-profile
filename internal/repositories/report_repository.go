package repositories

import (
	"log"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetRevenueReport(ownerID string, period string) (*models.RevenueReportResponse, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetRevenueReport(ownerID string, period string) (*models.RevenueReportResponse, error) {
	var response models.RevenueReportResponse
	response.Period = period

	log.Printf("ReportRepository: Getting revenue report for owner: %s, period: %s", ownerID, period)

	// Base query
	query := r.db.Debug().Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ?", ownerID, "paid")

	// Apply period filter
	switch period {
	case "day":
		query = query.Where("DATE(payments.paid_at) = CURRENT_DATE")
	case "month":
		query = query.Where("DATE_TRUNC('month', payments.paid_at) = DATE_TRUNC('month', CURRENT_DATE)")
	case "year":
		query = query.Where("DATE_TRUNC('year', payments.paid_at) = DATE_TRUNC('year', CURRENT_DATE)")
	}

	// Get total revenue for the period
	err := query.Session(&gorm.Session{}).Select("COALESCE(SUM(payments.amount), 0)").Scan(&response.TotalRevenue).Error
	if err != nil {
		log.Printf("ReportRepository: Total revenue error: %v", err)
		return nil, err
	}

	// Get breakdown by field
	err = query.Session(&gorm.Session{}).Select("fields.id as field_id, fields.name as field_name, COALESCE(SUM(payments.amount), 0) as revenue, COUNT(bookings.id) as booking_count").
		Group("fields.id, fields.name").
		Order("revenue DESC").
		Scan(&response.ByField).Error
	if err != nil {
		log.Printf("ReportRepository: Breakdown error: %v", err)
		return nil, err
	}

	// Get breakdown by sport type
	err = query.Session(&gorm.Session{}).Select("fields.sport_type, COALESCE(SUM(payments.amount), 0) as revenue, COUNT(bookings.id) as booking_count").
		Group("fields.sport_type").
		Order("revenue DESC").
		Scan(&response.BySportType).Error
	if err != nil {
		log.Printf("ReportRepository: Sport type breakdown error: %v", err)
		return nil, err
	}

	if response.ByField == nil {
		response.ByField = []models.RevenueByField{}
	}
	if response.BySportType == nil {
		response.BySportType = []models.RevenueBySportType{}
	}

	log.Printf("ReportRepository: Successfully fetched report. Total: %v, Field count: %v, Sport type count: %v", response.TotalRevenue, len(response.ByField), len(response.BySportType))

	return &response, nil
}
