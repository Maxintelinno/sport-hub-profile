package repositories

import (
	"time"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetActiveSubscription(userID string) (*models.Subscription, error)
	GetSummaryStats(userID string) (*models.DashboardSummary, error)
	GetRevenueTrend(userID string) ([]models.RevenueTrendItem, error)
	GetTodayBookingCount(userID string) (int64, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetActiveSubscription(userID string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Preload("Plan").
		Where("user_id = ? AND status IN ?", userID, []string{"trial", "active"}).
		Order("created_at DESC").
		First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *dashboardRepository) GetSummaryStats(userID string) (*models.DashboardSummary, error) {
	var stats models.DashboardSummary

	// Revenue & Booking Count
	err := r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ?", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0) as total_revenue, COUNT(DISTINCT bookings.id) as booking_count").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// Field Count
	r.db.Table("fields").Where("owner_id = ?", userID).Count(&stats.FieldCount)

	// Revenue Growth % (7d comparison)
	var current7d, previous7d float64
	r.db.Table("payments").
		Joins("JOIN bookings ON bookings.id = payments.booking_id").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Where("fields.owner_id = ? AND payments.status = 'paid' AND payments.paid_at >= CURRENT_DATE - INTERVAL '6 days'", userID).
		Select("COALESCE(SUM(payments.amount), 0)").Scan(&current7d)

	r.db.Table("payments").
		Joins("JOIN bookings ON bookings.id = payments.booking_id").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Where("fields.owner_id = ? AND payments.status = 'paid' AND payments.paid_at >= CURRENT_DATE - INTERVAL '13 days' AND payments.paid_at < CURRENT_DATE - INTERVAL '6 days'", userID).
		Select("COALESCE(SUM(payments.amount), 0)").Scan(&previous7d)

	if previous7d > 0 {
		stats.RevenueGrowthPct = ((current7d - previous7d) / previous7d) * 100
	} else if current7d > 0 {
		stats.RevenueGrowthPct = 100
	}

	return &stats, nil
}

func (r *dashboardRepository) GetRevenueTrend(userID string) ([]models.RevenueTrendItem, error) {
	var results []struct {
		Date   time.Time
		Amount float64
	}

	err := r.db.Raw(`
        SELECT 
            d.date,
            COALESCE(SUM(p.amount), 0) as amount
        FROM (
            SELECT generate_series(CURRENT_DATE - INTERVAL '6 days', CURRENT_DATE, '1 day')::date AS date
        ) d
        LEFT JOIN payments p ON DATE(p.paid_at) = d.date AND p.status = 'paid'
        LEFT JOIN bookings b ON b.id = p.booking_id
        LEFT JOIN fields f ON f.id = b.field_id AND f.owner_id = ?
        GROUP BY d.date
        ORDER BY d.date ASC
    `, userID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	trendItems := make([]models.RevenueTrendItem, len(results))
	thaiDays := map[string]string{
		"Monday": "จ.", "Tuesday": "อ.", "Wednesday": "พ.", "Thursday": "พฤ.", "Friday": "ศ.", "Saturday": "ส.", "Sunday": "อา.",
	}

	for i, res := range results {
		dayName := res.Date.Format("Monday")
		trendItems[i] = models.RevenueTrendItem{
			Label:  thaiDays[dayName],
			Amount: res.Amount,
		}
	}

	return trendItems, nil
}

func (r *dashboardRepository) GetTodayBookingCount(userID string) (int64, error) {
	var count int64
	err := r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Where("fields.owner_id = ? AND DATE(bookings.created_at) = CURRENT_DATE", userID).
		Count(&count).Error
	return count, err
}
