package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type ReportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetRevenueReport(c echo.Context) error {
	// Extract user ID from context (set by JWTMiddleware)
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
		})
	}

	period := c.QueryParam("group_by")
	// If the user specifies 'total', 'day', 'month', 'year', otherwise default to 'total'
	validPeriods := map[string]bool{"total": true, "day": true, "month": true, "year": true}
	if !validPeriods[period] {
		period = "total"
	}

	report, err := h.reportService.GetRevenueReport(userID, period)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, report)
}
