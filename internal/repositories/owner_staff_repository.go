package repositories

import (
	"log"
	"github.com/maxintelinno/sport-hub-profile/internal/models"

	"gorm.io/gorm"
)

type OwnerStaffRepository interface {
	GetStaffByOwnerID(ownerID string) ([]models.OwnerStaffResponse, error)
}

type ownerStaffRepository struct {
	db *gorm.DB
}

func NewOwnerStaffRepository(db *gorm.DB) OwnerStaffRepository {
	return &ownerStaffRepository{db: db}
}

func (r *ownerStaffRepository) GetStaffByOwnerID(ownerID string) ([]models.OwnerStaffResponse, error) {
	if !isUUID(ownerID) {
		return nil, gorm.ErrRecordNotFound
	}

	var staff []models.OwnerStaffResponse

	log.Printf("OwnerStaffRepository: Fetching staff for owner ID: %s", ownerID)

	err := r.db.Table("owner_staffs").
		Select("owner_staffs.id, owner_staffs.owner_user_id, owner_staffs.staff_user_id, users.username, users.fullname, users.phone, owner_staffs.role_code, owner_staffs.status, owner_staffs.created_at, owner_staffs.updated_at").
		Joins("JOIN users ON users.id = owner_staffs.staff_user_id").
		Where("owner_staffs.owner_user_id = ?", ownerID).
		Scan(&staff).Error

	if err != nil {
		log.Printf("OwnerStaffRepository: Error fetching staff for owner %s: %v", ownerID, err)
		return nil, err
	}

	return staff, nil
}
