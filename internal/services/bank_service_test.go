package services

import (
	"testing"

	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBankRepository struct {
	mock.Mock
}

func (m *MockBankRepository) GetByUserID(userID string) ([]models.OwnerBankAccount, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.OwnerBankAccount), args.Error(1)
}

func (m *MockBankRepository) GetByID(id string) (*models.OwnerBankAccount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OwnerBankAccount), args.Error(1)
}

func (m *MockBankRepository) Create(bankAccount *models.OwnerBankAccount) error {
	args := m.Called(bankAccount)
	return args.Error(0)
}

func (m *MockBankRepository) Update(bankAccount *models.OwnerBankAccount) error {
	args := m.Called(bankAccount)
	return args.Error(0)
}

func (m *MockBankRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBankRepository) UnsetDefaultByUserID(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockBankRepository) CountByUserID(userID string) (int64, error) {
	args := m.Called(userID)
	return int64(args.Int(0)), args.Error(1)
}

type MockPayoutRepository struct {
	mock.Mock
}

func (m *MockPayoutRepository) HasProcessingPayout(bankAccountID string) (bool, error) {
	args := m.Called(bankAccountID)
	return args.Bool(0), args.Error(1)
}

type BankMockUserRepository struct {
	mock.Mock
}

func (m *BankMockUserRepository) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Add other methods to BankMockUserRepository to satisfy interface
func (m *BankMockUserRepository) GetStatsByUserID(userID string) (*models.ProfileStats, error) { return nil, nil }
func (m *BankMockUserRepository) GetRevenueSummaryByUserID(userID string) (*models.RevenueSummary, error) { return nil, nil }
func (m *BankMockUserRepository) GetFieldCountByOwnerID(ownerID string) (int64, error) { return 0, nil }
func (m *BankMockUserRepository) GetCourtCountByOwnerID(ownerID string) (int64, error) { return 0, nil }
func (m *BankMockUserRepository) GetBookingCountByOwnerID(ownerID string) (int64, error) { return 0, nil }
func (m *BankMockUserRepository) GetUserByPhone(phone string) (*models.User, error) { return nil, nil }
func (m *BankMockUserRepository) UpdatePasswordByPhone(phone string, passwordHash string) error { return nil }
func (m *BankMockUserRepository) UpdatePinByPhone(phone string, pinHash string) error { return nil }
func (m *BankMockUserRepository) CleanupOTPs(phone string) error { return nil }
func (m *BankMockUserRepository) CreateOTP(otp *models.OTPRequest) error { return nil }

func TestMasking(t *testing.T) {
	t.Run("Mask Account Number", func(t *testing.T) {
		assert.Equal(t, "xxx-x-12345-x", maskAccountNumber("1234567890"))
		assert.Equal(t, "xxx-x-98765-x", maskAccountNumber("9876543210"))
	})

	t.Run("Mask PromptPay Phone", func(t *testing.T) {
		phone := "0903993838"
		masked := maskPromptpayValue(&phone)
		assert.Equal(t, "090-xxx-3838", *masked)
	})
}

func TestGetOwnerBankAccounts(t *testing.T) {
	mockBankRepo := new(MockBankRepository)
	mockUserRepo := new(BankMockUserRepository)
	mockPayoutRepo := new(MockPayoutRepository)
	s := NewBankService(mockBankRepo, mockUserRepo, mockPayoutRepo)

	t.Run("Success for owner", func(t *testing.T) {
		userID := "owner-1"
		mockUserRepo.On("GetUserByID", userID).Return(&models.User{ID: userID, Role: "เจ้าของสนาม"}, nil)
		mockBankRepo.On("GetByUserID", userID).Return([]models.OwnerBankAccount{
			{AccountNumber: "1234567890"},
		}, nil)

		resp, err := s.GetOwnerBankAccounts(userID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp))
		assert.Equal(t, "xxx-x-12345-x", resp[0].AccountNumberMasked)
	})

	t.Run("Error for non-owner", func(t *testing.T) {
		userID := "user-1"
		mockUserRepo.On("GetUserByID", userID).Return(&models.User{ID: userID, Role: "user"}, nil)

		resp, err := s.GetOwnerBankAccounts(userID)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

func TestUpdateOwnerBankAccount(t *testing.T) {
	mockBankRepo := new(MockBankRepository)
	mockUserRepo := new(BankMockUserRepository)
	mockPayoutRepo := new(MockPayoutRepository)
	s := NewBankService(mockBankRepo, mockUserRepo, mockPayoutRepo)

	userID := "owner-1"
	accountID := "acc-1"

	t.Run("Update sensitive info - resets verification", func(t *testing.T) {
		original := &models.OwnerBankAccount{ID: accountID, UserID: userID, AccountNumber: "111", IsVerified: true, VerificationStatus: "verified"}
		mockBankRepo.On("GetByID", accountID).Return(original, nil).Once()
		
		req := models.UpdateBankAccountRequest{AccountNumber: "222"}
		mockBankRepo.On("Update", mock.MatchedBy(func(a *models.OwnerBankAccount) bool {
			return a.AccountNumber == "222" && a.IsVerified == false && a.VerificationStatus == "pending"
		})).Return(nil).Once()

		err := s.UpdateOwnerBankAccount(userID, accountID, req)
		assert.NoError(t, err)
	})

	t.Run("Update non-sensitive info - keeps verification", func(t *testing.T) {
		original := &models.OwnerBankAccount{ID: accountID, UserID: userID, BankName: "Bank A", IsVerified: true, VerificationStatus: "verified"}
		mockBankRepo.On("GetByID", accountID).Return(original, nil).Once()
		
		req := models.UpdateBankAccountRequest{BankName: "Bank B"}
		mockBankRepo.On("Update", mock.MatchedBy(func(a *models.OwnerBankAccount) bool {
			return a.BankName == "Bank B" && a.IsVerified == true && a.VerificationStatus == "verified"
		})).Return(nil).Once()

		err := s.UpdateOwnerBankAccount(userID, accountID, req)
		assert.NoError(t, err)
	})
}

func TestSetDefaultBankAccount(t *testing.T) {
	mockBankRepo := new(MockBankRepository)
	mockUserRepo := new(BankMockUserRepository)
	mockPayoutRepo := new(MockPayoutRepository)
	s := NewBankService(mockBankRepo, mockUserRepo, mockPayoutRepo)

	userID := "owner-1"
	accountID := "acc-1"

	t.Run("Success", func(t *testing.T) {
		account := &models.OwnerBankAccount{ID: accountID, UserID: userID, IsDefault: false}
		mockBankRepo.On("GetByID", accountID).Return(account, nil).Once()
		mockBankRepo.On("UnsetDefaultByUserID", userID).Return(nil).Once()
		mockBankRepo.On("Update", mock.MatchedBy(func(a *models.OwnerBankAccount) bool {
			return a.ID == accountID && a.IsDefault == true
		})).Return(nil).Once()

		err := s.SetDefaultBankAccount(userID, accountID)
		assert.NoError(t, err)
	})
}

func TestDeleteOwnerBankAccount(t *testing.T) {
	mockBankRepo := new(MockBankRepository)
	mockUserRepo := new(BankMockUserRepository)
	mockPayoutRepo := new(MockPayoutRepository)
	s := NewBankService(mockBankRepo, mockUserRepo, mockPayoutRepo)

	userID := "owner-1"
	accountID := "acc-1"

	t.Run("Cannot delete only account", func(t *testing.T) {
		account := &models.OwnerBankAccount{ID: accountID, UserID: userID}
		mockBankRepo.On("GetByID", accountID).Return(account, nil).Once()
		mockBankRepo.On("CountByUserID", userID).Return(1, nil).Once()

		err := s.DeleteOwnerBankAccount(userID, accountID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only bank account")
	})

	t.Run("Cannot delete default with processing payouts", func(t *testing.T) {
		account := &models.OwnerBankAccount{ID: accountID, UserID: userID, IsDefault: true}
		mockBankRepo.On("GetByID", accountID).Return(account, nil).Once()
		mockPayoutRepo.On("HasProcessingPayout", accountID).Return(true, nil).Once()

		err := s.DeleteOwnerBankAccount(userID, accountID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payouts are processing")
	})

	t.Run("Success delete", func(t *testing.T) {
		account := &models.OwnerBankAccount{ID: accountID, UserID: userID, IsDefault: false}
		mockBankRepo.On("GetByID", accountID).Return(account, nil).Once()
		mockBankRepo.On("CountByUserID", userID).Return(2, nil).Once()
		mockBankRepo.On("Delete", accountID).Return(nil).Once()

		err := s.DeleteOwnerBankAccount(userID, accountID)
		assert.NoError(t, err)
	})
}
