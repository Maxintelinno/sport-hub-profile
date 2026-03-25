package repositories

import (
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"gorm.io/gorm"
)

type PayoutRepository interface {
	HasProcessingPayout(bankAccountID string) (bool, error)
}

type payoutRepository struct {
	db *gorm.DB
}

func NewPayoutRepository(db *gorm.DB) PayoutRepository {
	return &payoutRepository{db: db}
}

func (r *payoutRepository) HasProcessingPayout(bankAccountID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.OwnerPayout{}).
		Where("bank_account_id = ? AND status IN ?", bankAccountID, []string{"pending", "processing"}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
