package repositories

import (
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id string) (*models.User, error)
	GetStatsByUserID(userID string) (*models.ProfileStats, error)
	GetRevenueSummaryByUserID(userID string) (*models.RevenueSummary, error)
	GetFieldCountByOwnerID(ownerID string) (int64, error)
	GetCourtCountByOwnerID(ownerID string) (int64, error)
	GetBookingCountByOwnerID(ownerID string) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	// For now, return mock data since we might not have the table yet
	// In a real scenario, we would use: r.db.First(&user, id).Error
	user = models.User{
		ID:    id,
		Name:  "Owner Lastname",
		Phone: "0911113333",
		Role:  "เจ้าของสนาม",
	}
	return &user, nil
}

func (r *userRepository) GetFieldCountByOwnerID(ownerID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Field{}).Where("owner_id = ?", ownerID).Count(&count).Error
	return count, err
}

func (r *userRepository) GetCourtCountByOwnerID(ownerID string) (int64, error) {
	var count int64
	// Query to count courts belonging to any field owned by the owner
	err := r.db.Table("field_courts").
		Joins("JOIN fields ON fields.id = field_courts.field_id").
		Where("fields.owner_id = ?", ownerID).
		Count(&count).Error
	return count, err
}

func (r *userRepository) GetBookingCountByOwnerID(ownerID string) (int64, error) {
	var count int64
	// Query to count bookings for any field owned by the owner
	err := r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Where("fields.owner_id = ?", ownerID).
		Count(&count).Error
	return count, err
}

func (r *userRepository) GetStatsByUserID(userID string) (*models.ProfileStats, error) {
	fieldCount, err := r.GetFieldCountByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	bookingCount, err := r.GetBookingCountByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	// Assuming revenue still mocked for now
	return &models.ProfileStats{
		FieldCount:   int(fieldCount),
		BookingCount: int(bookingCount),
		TotalRevenue: 0,
	}, nil
}

func (r *userRepository) GetRevenueSummaryByUserID(userID string) (*models.RevenueSummary, error) {
	// Mock revenue summary
	return &models.RevenueSummary{
		Total:   0,
		Daily:   0,
		Monthly: 0,
		Yearly:  0,
	}, nil
}
