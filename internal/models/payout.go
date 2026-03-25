package models

import (
	"time"
)

type OwnerSettlement struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BookingID      string    `gorm:"type:uuid;not null;unique" json:"booking_id"`
	OwnerID        string    `gorm:"type:uuid;not null" json:"owner_id"`
	GrossAmount    float64   `gorm:"type:numeric(10,2);not null" json:"gross_amount"`
	PlatformFee    float64   `gorm:"type:numeric(10,2);not null;default:0" json:"platform_fee"`
	DiscountAmount float64   `gorm:"type:numeric(10,2);not null;default:0" json:"discount_amount"`
	NetAmount      float64   `gorm:"type:numeric(10,2);not null" json:"net_amount"`
	Status         string    `gorm:"type:varchar(30);not null;default:'pending'" json:"status"`
	AvailableAt    *time.Time `json:"available_at"`
	PaidAt         *time.Time `json:"paid_at"`
	CreatedAt      time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

type OwnerPayout struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PayoutNo          string    `gorm:"type:varchar(30);not null;unique" json:"payout_no"`
	OwnerID           string    `gorm:"type:uuid;not null" json:"owner_id"`
	BankAccountID     string    `gorm:"type:uuid;not null" json:"bank_account_id"`
	TotalAmount       float64   `gorm:"type:numeric(10,2);not null" json:"total_amount"`
	TransferFee       float64   `gorm:"type:numeric(10,2);not null;default:0" json:"transfer_fee"`
	FinalAmount       float64   `gorm:"type:numeric(10,2);not null" json:"final_amount"`
	Status            string    `gorm:"type:varchar(30);not null;default:'pending'" json:"status"`
	PayoutMethod      string    `gorm:"type:varchar(30);not null;default:'bank_transfer'" json:"payout_method"`
	RequestedAt       *time.Time `json:"requested_at"`
	ProcessedAt       *time.Time `json:"processed_at"`
	PaidAt            *time.Time `json:"paid_at"`
	FailedAt          *time.Time `json:"failed_at"`
	CancelledAt       *time.Time `json:"cancelled_at"`
	FailureReason     *string   `json:"failure_reason"`
	TransferReference *string   `json:"transfer_reference"`
	TransferProvider  *string   `json:"transfer_provider"`
	CreatedAt         time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}
