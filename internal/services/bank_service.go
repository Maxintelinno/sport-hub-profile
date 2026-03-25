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
	UpdateOwnerBankAccount(userID string, accountID string, req models.UpdateBankAccountRequest) error
	SetDefaultBankAccount(userID string, accountID string) error
	DeleteOwnerBankAccount(userID string, accountID string) error
}

type bankService struct {
	bankRepo   repositories.BankRepository
	userRepo   repositories.UserRepository
	payoutRepo repositories.PayoutRepository
}

func NewBankService(bankRepo repositories.BankRepository, userRepo repositories.UserRepository, payoutRepo repositories.PayoutRepository) BankService {
	return &bankService{bankRepo: bankRepo, userRepo: userRepo, payoutRepo: payoutRepo}
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

func (s *bankService) UpdateOwnerBankAccount(userID string, accountID string, req models.UpdateBankAccountRequest) error {
	// 1. Get original account
	account, err := s.bankRepo.GetByID(accountID)
	if err != nil {
		return err
	}

	// 2. Verify ownership
	if account.UserID != userID {
		return errors.New("unauthorized: account does not belong to you")
	}

	// 3. Check for sensitive changes
	sensitiveChanged := false
	if req.AccountNumber != "" && req.AccountNumber != account.AccountNumber {
		sensitiveChanged = true
		account.AccountNumber = req.AccountNumber
	}
	if req.PromptpayValue != nil && (account.PromptpayValue == nil || *req.PromptpayValue != *account.PromptpayValue) {
		sensitiveChanged = true
		account.PromptpayValue = req.PromptpayValue
	}

	if sensitiveChanged {
		account.IsVerified = false
		account.VerificationStatus = "pending"
	}

	// 4. Update other fields
	if req.BankCode != "" {
		account.BankCode = req.BankCode
	}
	if req.BankName != "" {
		account.BankName = req.BankName
	}
	if req.AccountName != "" {
		account.AccountName = req.AccountName
	}
	if req.PromptpayType != nil {
		account.PromptpayType = req.PromptpayType
	}

	// 5. Handle default logic
	if req.IsDefault && !account.IsDefault {
		err = s.bankRepo.UnsetDefaultByUserID(userID)
		if err != nil {
			log.Printf("BankService: Error unsetting defaults: %v", err)
		}
		account.IsDefault = true
	}

	return s.bankRepo.Update(account)
}

func (s *bankService) SetDefaultBankAccount(userID string, accountID string) error {
	// 1. Get account
	account, err := s.bankRepo.GetByID(accountID)
	if err != nil {
		return err
	}

	// 2. Verify ownership
	if account.UserID != userID {
		return errors.New("unauthorized: account does not belong to you")
	}

	// 3. Unset Others
	err = s.bankRepo.UnsetDefaultByUserID(userID)
	if err != nil {
		return err
	}

	// 4. Set this one
	account.IsDefault = true
	return s.bankRepo.Update(account)
}

func (s *bankService) DeleteOwnerBankAccount(userID string, accountID string) error {
	// 1. Get account
	account, err := s.bankRepo.GetByID(accountID)
	if err != nil {
		return err
	}

	// 2. Verify ownership
	if account.UserID != userID {
		return errors.New("unauthorized: account does not belong to you")
	}

	// 3. Check if default and has processing payouts
	if account.IsDefault {
		hasProcessing, err := s.payoutRepo.HasProcessingPayout(accountID)
		if err != nil {
			return err
		}
		if hasProcessing {
			return errors.New("cannot delete default account while payouts are processing")
		}
	}

	// 4. Check if only account
	count, err := s.bankRepo.CountByUserID(userID)
	if err != nil {
		return err
	}
	if count <= 1 {
		return errors.New("cannot delete the only bank account. please add another one first.")
	}

	return s.bankRepo.Delete(accountID)
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
