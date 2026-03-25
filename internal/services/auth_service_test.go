package services

import (
	"errors"
	"testing"

	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetStatsByUserID(userID string) (*models.ProfileStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProfileStats), args.Error(1)
}

func (m *MockUserRepository) GetRevenueSummaryByUserID(userID string) (*models.RevenueSummary, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RevenueSummary), args.Error(1)
}

func (m *MockUserRepository) GetFieldCountByOwnerID(ownerID string) (int64, error) {
	args := m.Called(ownerID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockUserRepository) GetCourtCountByOwnerID(ownerID string) (int64, error) {
	args := m.Called(ownerID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockUserRepository) GetBookingCountByOwnerID(ownerID string) (int64, error) {
	args := m.Called(ownerID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockUserRepository) GetUserByPhone(phone string) (*models.User, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePasswordByPhone(phone string, passwordHash string) error {
	args := m.Called(phone, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) CleanupOTPs(phone string) error {
	args := m.Called(phone)
	return args.Error(0)
}

func (m *MockUserRepository) CreateOTP(otp *models.OTPRequest) error {
	args := m.Called(otp)
	return args.Error(0)
}

func TestCheckPhoneWithOTP(t *testing.T) {
	t.Run("Phone exists - triggers OTP process", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		phone := "0123456789"
		user := &models.User{Phone: phone}

		mockRepo.On("GetUserByPhone", phone).Return(user, nil)
		mockRepo.On("CleanupOTPs", phone).Return(nil)
		mockRepo.On("CreateOTP", mock.AnythingOfType("*models.OTPRequest")).Return(nil)

		s := NewAuthService(mockRepo)
		exists, err := s.CheckPhone(phone)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Phone does not exist", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		phone := "9999999999"

		mockRepo.On("GetUserByPhone", phone).Return(nil, errors.New("not found"))

		s := NewAuthService(mockRepo)
		exists, err := s.CheckPhone(phone)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertNotCalled(t, "CleanupOTPs", mock.Anything)
		mockRepo.AssertNotCalled(t, "CreateOTP", mock.Anything)
	})
}
