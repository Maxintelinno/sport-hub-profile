package repositories

import (
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint) (*models.User, error)
	GetStatsByUserID(userID uint) (*models.ProfileStats, error)
	GetRevenueSummaryByUserID(userID uint) (*models.RevenueSummary, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
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

func (r *userRepository) GetStatsByUserID(userID uint) (*models.ProfileStats, error) {
	// Mock stats
	return &models.ProfileStats{
		FieldCount:   0,
		BookingCount: 0,
		TotalRevenue: 0,
	}, nil
}

func (r *userRepository) GetRevenueSummaryByUserID(userID uint) (*models.RevenueSummary, error) {
	// Mock revenue summary
	return &models.RevenueSummary{
		Total:   0,
		Daily:   0,
		Monthly: 0,
		Yearly:  0,
	}, nil
}
