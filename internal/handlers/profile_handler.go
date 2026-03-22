package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type ProfileHandler struct {
	profileService services.ProfileService
}

func NewProfileHandler(profileService services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(c echo.Context) error {
	// For now, assume user ID 1
	userID := uint(1)
	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, profile)
}
