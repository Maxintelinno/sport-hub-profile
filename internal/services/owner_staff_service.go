package services

import (
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
)

type OwnerStaffService interface {
	GetStaffByOwnerID(ownerID string) ([]models.OwnerStaffResponse, error)
}

type ownerStaffService struct {
	repo repositories.OwnerStaffRepository
}

func NewOwnerStaffService(repo repositories.OwnerStaffRepository) OwnerStaffService {
	return &ownerStaffService{repo: repo}
}

func (s *ownerStaffService) GetStaffByOwnerID(ownerID string) ([]models.OwnerStaffResponse, error) {
	return s.repo.GetStaffByOwnerID(ownerID)
}
