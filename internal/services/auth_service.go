package services

import (
	"errors"
	"log"

	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ForgotPassword(phone string, newPassword string) error
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
