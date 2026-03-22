package repositories

import (
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

	// Base query
	query := r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Where("fields.owner_id = ? AND bookings.payment_status = ?", ownerID, "paid")

	// Apply period filter for total and breakdown
	switch period {
	case "day":
		query = query.Where("DATE(bookings.paid_at) = CURRENT_DATE")
	case "month":
		query = query.Where("DATE_TRUNC('month', bookings.paid_at) = DATE_TRUNC('month', CURRENT_DATE)")
	case "year":
		query = query.Where("DATE_TRUNC('year', bookings.paid_at) = DATE_TRUNC('year', CURRENT_DATE)")
	// "total" case doesn't need additional where
	}

	// Get total revenue for the period
	err := query.Session(&gorm.Session{}).Select("COALESCE(SUM(bookings.total_amount), 0)").Scan(&response.TotalRevenue).Error
	if err != nil {
		return nil, err
	}

	// Get breakdown by field
	err = query.Session(&gorm.Session{}).Select("fields.id as field_id, fields.name as field_name, COALESCE(SUM(bookings.total_amount), 0) as revenue, COUNT(bookings.id) as booking_count").
		Group("fields.id, fields.name").
		Order("revenue DESC").
		Scan(&response.ByField).Error
	if err != nil {
		return nil, err
	}

	if response.ByField == nil {
		response.ByField = []models.RevenueByField{}
	}

	return &response, nil
}
