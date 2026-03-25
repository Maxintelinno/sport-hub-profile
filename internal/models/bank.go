package models

import "time"

type OwnerBankAccount struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID             string    `gorm:"type:uuid;not null" json:"user_id"`
	BankCode           string    `gorm:"type:varchar(20);not null" json:"bank_code"`
	BankName           string    `gorm:"type:varchar(100);not null" json:"bank_name"`
	AccountName        string    `gorm:"type:varchar(255);not null" json:"account_name"`
	AccountNumber      string    `gorm:"type:varchar(50);not null" json:"account_number"`
	PromptpayType      *string   `gorm:"type:varchar(20)" json:"promptpay_type"`
	PromptpayValue     *string   `gorm:"type:varchar(100)" json:"promptpay_value"`
	IsDefault          bool      `gorm:"type:boolean;not null;default:false" json:"is_default"`
	IsVerified         bool      `gorm:"type:boolean;not null;default:false" json:"is_verified"`
	VerificationStatus string    `gorm:"type:varchar(30);not null;default:'pending'" json:"verification_status"`
	Status             string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	Note               *string   `gorm:"type:text" json:"note"`
	CreatedAt          time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

type AddBankAccountRequest struct {
	BankCode       string  `json:"bank_code"`
	BankName       string  `json:"bank_name"`
	AccountName    string  `json:"account_name" validate:"required"`
	AccountNumber  string  `json:"account_number"`
	PromptpayType  *string `json:"promptpay_type"`
	PromptpayValue *string `json:"promptpay_value"`
	IsDefault      bool    `json:"is_default"`
}

type BankAccountResponse struct {
	ID                   string    `json:"id"`
	BankCode             string    `json:"bank_code"`
	BankName             string    `json:"bank_name"`
	AccountName          string    `json:"account_name"`
	AccountNumberMasked  string    `json:"account_number_masked"`
	PromptpayType        *string   `json:"promptpay_type"`
	PromptpayValueMasked *string   `json:"promptpay_value_masked"`
	IsDefault            bool      `json:"is_default"`
	IsVerified           bool      `json:"is_verified"`
	VerificationStatus   string    `json:"verification_status"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type BankAccountListResponse struct {
	Status         string                `json:"status"`
	Message        string                `json:"message"`
	HasBankAccount bool                  `json:"has_bank_account"`
	Data           []BankAccountResponse `json:"data"`
}

type BankAccountCreateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
}
