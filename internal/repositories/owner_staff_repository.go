package repositories

import (
	"log"
	"github.com/maxintelinno/sport-hub-profile/internal/models"

	"gorm.io/gorm"
)

type OwnerStaffRepository interface {
	GetStaffByOwnerID(ownerID string) ([]models.OwnerStaffResponse, error)
	UpdateStaffStatus(ownerID string, staffUserID string, status string) error
	GetOwnerIDByStaffUserID(staffUserID string) (string, error)
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

func (r *ownerStaffRepository) UpdateStaffStatus(ownerID string, staffUserID string, status string) error {
	if !isUUID(ownerID) || !isUUID(staffUserID) {
		return gorm.ErrRecordNotFound
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update owner_staffs first to check if owner owns this staff
		log.Printf("OwnerStaffRepository: Updating status to %s for staff %s by owner %s", status, staffUserID, ownerID)
		
		result := tx.Model(&models.OwnerStaff{}).
			Where("owner_user_id = ? AND staff_user_id = ?", ownerID, staffUserID).
			Update("status", status)
		
		if result.Error != nil {
			return result.Error
		}
		
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		// Update users
		if err := tx.Model(&models.User{}).
			Where("id = ?", staffUserID).
			Update("status", status).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ownerStaffRepository) GetOwnerIDByStaffUserID(staffUserID string) (string, error) {
	if !isUUID(staffUserID) {
		return "", gorm.ErrRecordNotFound
	}

	var ownerStaff models.OwnerStaff
	log.Printf("OwnerStaffRepository: Finding owner ID for staff user ID: %s", staffUserID)
	
	err := r.db.Where("staff_user_id = ? AND status = ?", staffUserID, "active").First(&ownerStaff).Error
	if err != nil {
		log.Printf("OwnerStaffRepository: Error finding owner for staff %s: %v", staffUserID, err)
		return "", err
	}

	return ownerStaff.OwnerUserID, nil
}
