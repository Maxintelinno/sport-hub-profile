package services

import (
	"strings"

	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
)

type ProfileService interface {
	GetProfile(userID uint) (*models.ProfileResponse, error)
}

type profileService struct {
	userRepo repositories.UserRepository
}

func NewProfileService(userRepo repositories.UserRepository) ProfileService {
	return &profileService{userRepo: userRepo}
}

func (s *profileService) GetProfile(userID uint) (*models.ProfileResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	stats, err := s.userRepo.GetStatsByUserID(userID)
	if err != nil {
		return nil, err
	}

	revenue, err := s.userRepo.GetRevenueSummaryByUserID(userID)
	if err != nil {
		return nil, err
	}

	initials := ""
	if len(user.Name) > 0 {
		initials = strings.ToUpper(string(user.Name[0]))
	}

	return &models.ProfileResponse{
		User: models.UserSummary{
			Name:      user.Name,
			Phone:     user.Phone,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
			Initials:  initials,
		},
		Stats: *stats,
		Plan: models.PlanInfo{
			Name:       "Free Plan",
			FieldUsage: "1/1",
			CourtUsage: "2/2",
			CanUpgrade: true,
		},
		RevenueSummary: *revenue,
	}, nil
}
