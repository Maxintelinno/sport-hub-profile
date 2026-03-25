package repositories

import (
	"log"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"gorm.io/gorm"
)

type BankRepository interface {
	GetByUserID(userID string) ([]models.OwnerBankAccount, error)
	Create(bankAccount *models.OwnerBankAccount) error
	UnsetDefaultByUserID(userID string) error
}

type bankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &bankRepository{db: db}
}

func (r *bankRepository) GetByUserID(userID string) ([]models.OwnerBankAccount, error) {
	var accounts []models.OwnerBankAccount
	log.Printf("BankRepository: Fetching accounts for user: %s", userID)
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		log.Printf("BankRepository: Error fetching accounts: %v", err)
		return nil, err
	}
	return accounts, nil
}

func (r *bankRepository) Create(bankAccount *models.OwnerBankAccount) error {
	log.Printf("BankRepository: Creating account for user: %s", bankAccount.UserID)
	return r.db.Create(bankAccount).Error
}

func (r *bankRepository) UnsetDefaultByUserID(userID string) error {
	log.Printf("BankRepository: Unsetting default accounts for user: %s", userID)
	return r.db.Model(&models.OwnerBankAccount{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}
