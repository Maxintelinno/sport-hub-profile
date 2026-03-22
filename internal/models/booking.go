package models

import (
	"time"
)

type Booking struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BookingNo     string    `gorm:"type:varchar(30);unique;not null" json:"booking_no"`
	UserID        string    `gorm:"type:uuid;not null" json:"user_id"`
	FieldID       string    `gorm:"type:uuid;not null" json:"field_id"`
	BookingDate   time.Time `gorm:"type:date;not null" json:"booking_date"`
	TotalAmount   float64   `gorm:"type:numeric(10,2);default:0;not null" json:"total_amount"`
	Status        string    `gorm:"type:varchar(20);default:'pending';not null" json:"status"`
	PaymentStatus string    `gorm:"type:varchar(20);default:'unpaid';not null" json:"payment_status"`
	Note          string    `gorm:"type:text" json:"note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	PaidAt        *time.Time `json:"paid_at"`
	CancelledAt   *time.Time `json:"cancelled_at"`
	CancelReason  string    `gorm:"type:text" json:"cancel_reason"`
}
