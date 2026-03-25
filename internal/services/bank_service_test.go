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

func (m *MockBankRepository) Create(bankAccount *models.OwnerBankAccount) error {
	args := m.Called(bankAccount)
	return args.Error(0)
}

func (m *MockBankRepository) UnsetDefaultByUserID(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
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
	s := NewBankService(mockBankRepo, mockUserRepo)

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
