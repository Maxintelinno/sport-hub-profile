package services

import (
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
)

type ReportService interface {
	GetRevenueReport(userID string, period string) (*models.RevenueReportResponse, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GetRevenueReport(userID string, period string) (*models.RevenueReportResponse, error) {
	if period == "" {
		period = "total"
	}
	return s.reportRepo.GetRevenueReport(userID, period)
}
