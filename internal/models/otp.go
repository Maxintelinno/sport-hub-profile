package models

import "time"

type OTPRequest struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Phone     string    `gorm:"type:varchar(20);not null" json:"phone"`
	OTPHash   string    `gorm:"type:text;not null" json:"otp_hash"`
	Attempts  int       `gorm:"type:int4;default:0" json:"attempts"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
}

func (OTPRequest) TableName() string {
	return "otp_requests"
}
