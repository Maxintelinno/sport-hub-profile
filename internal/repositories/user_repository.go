package repositories

import (
	"log"
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
	GetUserByPhone(phone string) (*models.User, error)
	UpdatePasswordByPhone(phone string, passwordHash string) error
	UpdatePinByPhone(phone string, pinHash string) error
	CleanupOTPs(phone string) error
	CreateOTP(otp *models.OTPRequest) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	if !isUUID(id) {
		return nil, gorm.ErrRecordNotFound
	}
	var user models.User
	log.Printf("UserRepository: Fetching user by ID: %s", id)
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		log.Printf("UserRepository: Error fetching user %s: %v", id, err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetFieldCountByOwnerID(ownerID string) (int64, error) {
	if !isUUID(ownerID) {
		return 0, nil
	}
	var count int64
	log.Printf("UserRepository: Getting field count for owner: %s", ownerID)
	err := r.db.Model(&models.Field{}).Where("owner_id = ?", ownerID).Count(&count).Error
	if err != nil {
		log.Printf("UserRepository: Error counting fields: %v", err)
	}
	log.Printf("UserRepository: Field count result: %d", count)
	return count, err
}

func (r *userRepository) GetCourtCountByOwnerID(ownerID string) (int64, error) {
	if !isUUID(ownerID) {
		return 0, nil
	}
	var count int64
	log.Printf("UserRepository: Getting court count for owner: %s", ownerID)
	// Query to count courts belonging to any field owned by the owner
	err := r.db.Table("field_courts").
		Joins("JOIN fields ON fields.id = field_courts.field_id").
		Where("fields.owner_id = ?", ownerID).
		Count(&count).Error
	if err != nil {
		log.Printf("UserRepository: Error counting courts: %v", err)
	}
	log.Printf("UserRepository: Court count result: %d", count)
	return count, err
}

func (r *userRepository) GetBookingCountByOwnerID(ownerID string) (int64, error) {
	if !isUUID(ownerID) {
		return 0, nil
	}
	var count int64
	log.Printf("UserRepository: Getting booking count for owner: %s", ownerID)
	// Query to count bookings for any field owned by the owner
	err := r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Where("fields.owner_id = ?", ownerID).
		Count(&count).Error
	if err != nil {
		log.Printf("UserRepository: Error counting bookings: %v", err)
	}
	log.Printf("UserRepository: Booking count result: %d", count)
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

	var totalRevenue float64
	err = r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ?", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0)").
		Scan(&totalRevenue).Error
	if err != nil {
		log.Printf("UserRepository: Error calculating total revenue: %v", err)
	}

	return &models.ProfileStats{
		FieldCount:   int(fieldCount),
		BookingCount: int(bookingCount),
		TotalRevenue: totalRevenue,
	}, nil
}

func (r *userRepository) GetRevenueSummaryByUserID(userID string) (*models.RevenueSummary, error) {
	var summary models.RevenueSummary
	if !isUUID(userID) {
		return &summary, nil
	}

	// Total
	err := r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ?", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0)").
		Scan(&summary.Total).Error
	if err != nil {
		return nil, err
	}

	// Daily
	err = r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ? AND DATE(payments.paid_at) = CURRENT_DATE", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0)").
		Scan(&summary.Daily).Error
	if err != nil {
		return nil, err
	}

	// Monthly
	err = r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ? AND DATE_TRUNC('month', payments.paid_at) = DATE_TRUNC('month', CURRENT_DATE)", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0)").
		Scan(&summary.Monthly).Error
	if err != nil {
		return nil, err
	}

	// Yearly
	err = r.db.Table("bookings").
		Joins("JOIN fields ON fields.id = bookings.field_id").
		Joins("JOIN payments ON payments.booking_id = bookings.id").
		Where("fields.owner_id = ? AND payments.status = ? AND DATE_TRUNC('year', payments.paid_at) = DATE_TRUNC('year', CURRENT_DATE)", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0)").
		Scan(&summary.Yearly).Error
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
func (r *userRepository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	log.Printf("UserRepository: Fetching user by phone: %s", phone)
	if err := r.db.First(&user, "phone = ?", phone).Error; err != nil {
		log.Printf("UserRepository: Error fetching user by phone %s: %v", phone, err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePasswordByPhone(phone string, passwordHash string) error {
	log.Printf("UserRepository: Updating password for phone: %s", phone)
	result := r.db.Model(&models.User{}).Where("phone = ?", phone).Update("password_hash", passwordHash)
	if result.Error != nil {
		log.Printf("UserRepository: Error updating password for phone %s: %v", phone, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) UpdatePinByPhone(phone string, pinHash string) error {
	log.Printf("UserRepository: Updating PIN for phone: %s", phone)
	result := r.db.Model(&models.User{}).Where("phone = ?", phone).Update("pin_hash", pinHash)
	if result.Error != nil {
		log.Printf("UserRepository: Error updating PIN for phone %s: %v", phone, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) CleanupOTPs(phone string) error {
	log.Printf("UserRepository: Cleaning up old OTPs for phone: %s", phone)
	err := r.db.Where("phone = ?", phone).Delete(&models.OTPRequest{}).Error
	if err != nil {
		log.Printf("UserRepository: Error cleaning up OTPs for phone %s: %v", phone, err)
	}
	return err
}

func (r *userRepository) CreateOTP(otp *models.OTPRequest) error {
	log.Printf("UserRepository: Creating new OTP for phone: %s", otp.Phone)
	err := r.db.Create(otp).Error
	if err != nil {
		log.Printf("UserRepository: Error creating OTP for phone %s: %v", otp.Phone, err)
	}
	return err
}
