package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
)

type BankService interface {
	GetOwnerBankAccounts(userID string) ([]models.BankAccountResponse, error)
	AddOwnerBankAccount(userID string, req models.AddBankAccountRequest) (*models.OwnerBankAccount, error)
}

type bankService struct {
	bankRepo repositories.BankRepository
	userRepo repositories.UserRepository
}

func NewBankService(bankRepo repositories.BankRepository, userRepo repositories.UserRepository) BankService {
	return &bankService{bankRepo: bankRepo, userRepo: userRepo}
}

func (s *bankService) GetOwnerBankAccounts(userID string) ([]models.BankAccountResponse, error) {
	// 1. Verify user exists and is an owner
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// For simplicity, checking if role contains "เจ้าของ" or "owner"
	// As seen in tests, role is "เจ้าของสนาม"
	if !strings.Contains(user.Role, "เจ้าของ") && !strings.Contains(strings.ToLower(user.Role), "owner") {
		return nil, errors.New("unauthorized: only owners can access bank accounts")
	}

	accounts, err := s.bankRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var resp []models.BankAccountResponse
	for _, acc := range accounts {
		resp = append(resp, models.BankAccountResponse{
			ID:                   acc.ID,
			BankCode:             acc.BankCode,
			BankName:             acc.BankName,
			AccountName:          acc.AccountName,
			AccountNumberMasked:  maskAccountNumber(acc.AccountNumber),
			PromptpayType:        acc.PromptpayType,
			PromptpayValueMasked: maskPromptpayValue(acc.PromptpayValue),
			IsDefault:            acc.IsDefault,
			IsVerified:           acc.IsVerified,
			VerificationStatus:   acc.VerificationStatus,
			Status:               acc.Status,
			CreatedAt:            acc.CreatedAt,
			UpdatedAt:            acc.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *bankService) AddOwnerBankAccount(userID string, req models.AddBankAccountRequest) (*models.OwnerBankAccount, error) {
	// 1. Verify user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(user.Role, "เจ้าของ") && !strings.Contains(strings.ToLower(user.Role), "owner") {
		return nil, errors.New("unauthorized: only owners can add bank accounts")
	}

	// 2. Validate request
	if req.AccountName == "" {
		return nil, errors.New("account_name is required")
	}

	hasBank := req.BankCode != "" && req.AccountNumber != ""
	hasPromptpay := req.PromptpayType != nil && req.PromptpayValue != nil && *req.PromptpayValue != ""

	if !hasBank && !hasPromptpay {
		return nil, errors.New("either bank account or promptpay info must be provided")
	}

	// 3. Handle default logic
	if req.IsDefault {
		err = s.bankRepo.UnsetDefaultByUserID(userID)
		if err != nil {
			log.Printf("BankService: Error unsetting defaults: %v", err)
		}
	}

	// 4. Create account
	bankAccount := &models.OwnerBankAccount{
		UserID:         userID,
		BankCode:       req.BankCode,
		BankName:       req.BankName,
		AccountName:    req.AccountName,
		AccountNumber:  req.AccountNumber,
		PromptpayType:  req.PromptpayType,
		PromptpayValue: req.PromptpayValue,
		IsDefault:      req.IsDefault,
		Status:         "active",
		VerificationStatus: "pending",
	}

	err = s.bankRepo.Create(bankAccount)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func maskAccountNumber(num string) string {
	if len(num) < 6 {
		return num
	}
	// Example: 1234567890 -> xxx-x-12345-x
	// Following user's example pattern: mask some parts, show middle?
	// The user's example shows 12345 which are the FIRST 5 digits of 1234567890.
	// Format: xxx-x-[first 5]-x
	if len(num) >= 5 {
		return fmt.Sprintf("xxx-x-%s-x", num[:5])
	}
	return "xxx-x-xxxx-x"
}

func maskPromptpayValue(val *string) *string {
	if val == nil || *val == "" {
		return nil
	}
	v := *val
	if len(v) < 6 {
		masked := "xxx-xxx-xxxx"
		return &masked
	}
	// Simple mask for phone or ID
	masked := ""
	if len(v) == 10 { // Phone
		masked = fmt.Sprintf("%s-xxx-%s", v[:3], v[6:])
	} else if len(v) == 13 { // National ID
		masked = fmt.Sprintf("%s-xxxx-xxxx-%s", v[:1], v[12:])
	} else {
		masked = "xxx-xxxx-" + v[len(v)-4:]
	}
	return &masked
}
