package services

import (
	"fmt"
	"strings"

	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
)

type ProfileService interface {
	GetProfile(userID string) (*models.ProfileResponse, error)
}

type profileService struct {
	userRepo repositories.UserRepository
}

func NewProfileService(userRepo repositories.UserRepository) ProfileService {
	return &profileService{userRepo: userRepo}
}

func (s *profileService) GetProfile(userID string) (*models.ProfileResponse, error) {
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

	courtCount, err := s.userRepo.GetCourtCountByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	initials := ""
	if len(user.Fullname) > 0 {
		initials = strings.ToUpper(string([]rune(user.Fullname)[0]))
	}

	return &models.ProfileResponse{
		User: models.UserSummary{
			Name:      user.Fullname,
			Phone:     user.Phone,
			Role:      user.Role,
			AvatarURL: "", // AvatarURL not in schema
			Initials:  initials,
		},
		Stats: *stats,
		Plan: models.PlanInfo{
			Name:       "Free Plan",
			FieldUsage: fmt.Sprintf("%d/1", stats.FieldCount),
			CourtUsage: fmt.Sprintf("%d/2", courtCount),
			CanUpgrade: true,
		},
		RevenueSummary: *revenue,
	}, nil
}
