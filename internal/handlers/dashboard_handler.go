package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type DashboardHandler interface {
	GetDashboard(c echo.Context) error
}

type dashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) DashboardHandler {
	return &dashboardHandler{dashboardService: dashboardService}
}

func (h *dashboardHandler) GetDashboard(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Unauthorized or invalid user ID"})
	}

	response, err := h.dashboardService.GetDashboard(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Owner dashboard retrieved successfully",
		"data":    response,
	})
}
