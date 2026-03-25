package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type AuthHandler interface {
	ForgotPassword(c echo.Context) error
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) ForgotPassword(c echo.Context) error {
	var req models.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Invalid request body"})
	}

	if req.Phone == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Phone and new_password are required"})
	}

	if len(req.NewPassword) < 6 {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Password must be at least 6 characters"})
	}

	err := h.authService.ForgotPassword(req.Phone, req.NewPassword)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.AuthResponse{Message: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.AuthResponse{Message: "Error updating password"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{Message: "Password updated successfully"})
}
