package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ForgotPassword(phone string, newPassword string) error
	CheckPhone(phone string) (bool, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) ForgotPassword(phone string, newPassword string) error {
	log.Printf("AuthService: Starting forgot password for phone: %s", phone)

	// 1. Check if user exists
	_, err := s.userRepo.GetUserByPhone(phone)
	if err != nil {
		log.Printf("AuthService: User not found for phone %s: %v", phone, err)
		return errors.New("user not found")
	}

	// 2. Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("AuthService: Error hashing password: %v", err)
		return err
	}

	// 3. Update password in database
	err = s.userRepo.UpdatePasswordByPhone(phone, string(hashedPassword))
	if err != nil {
		log.Printf("AuthService: Error updating password: %v", err)
		return err
	}

	log.Printf("AuthService: Password updated successfully for phone: %s", phone)
	return nil
}

func (s *authService) CheckPhone(phone string) (bool, error) {
	log.Printf("AuthService: Checking if phone exists: %s", phone)
	_, err := s.userRepo.GetUserByPhone(phone)
	if err != nil {
		log.Printf("AuthService: Phone %s not found: %v", phone, err)
		return false, nil // Phone not found is not an error here, just return false
	}

	// Phone exists, proceed with OTP logic
	log.Printf("AuthService: Phone %s exists, initiating OTP process", phone)

	// 1. Clean up old OTPs for this phone
	err = s.userRepo.CleanupOTPs(phone)
	if err != nil {
		log.Printf("AuthService: Error cleaning up old OTPs for %s: %v", phone, err)
		// We can continue even if cleanup fails, but log it
	}

	// 2. Generate 6-digit OTP
	otp, err := s.generateOTP()
	if err != nil {
		log.Printf("AuthService: Error generating OTP for %s: %v", phone, err)
		return true, err
	}
	log.Printf("AuthService: Generated OTP for %s (for dev/log): %s", phone, otp)

	// 3. Hash OTP
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("AuthService: Error hashing OTP for %s: %v", phone, err)
		return true, err
	}

	// 4. Save OTP request
	otpRequest := &models.OTPRequest{
		Phone:     phone,
		OTPHash:   string(otpHash),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	err = s.userRepo.CreateOTP(otpRequest)
	if err != nil {
		log.Printf("AuthService: Error saving OTP request for %s: %v", phone, err)
		return true, err
	}

	return true, nil
}

func (s *authService) generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}
