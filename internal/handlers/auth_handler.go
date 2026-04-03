package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type AuthHandler interface {
	ForgotPassword(c echo.Context) error
	CheckPhone(c echo.Context) error
	UpdatePassword(c echo.Context) error
	UpdatePin(c echo.Context) error
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

func (h *authHandler) CheckPhone(c echo.Context) error {
	var req models.CheckPhoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Invalid request body"})
	}

	if req.Phone == "" {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Phone is required"})
	}

	exists, err := h.authService.CheckPhone(req.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.CheckPhoneResponse{
			Message: "Error checking phone",
			IsFound: false,
		})
	}

	if !exists {
		return c.JSON(http.StatusNotFound, models.CheckPhoneResponse{
			Message: "ไม่มีเบอร์โทรนี้ในระบบ กรุณาตรวจสอบใหม่อีกครั้ง",
			IsFound: false,
		})
	}

	return c.JSON(http.StatusOK, models.CheckPhoneResponse{
		Message: "Phone number is registered",
		IsFound: true,
	})
}
func (h *authHandler) UpdatePassword(c echo.Context) error {
	var req models.UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Invalid request body"})
	}

	if req.Phone == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Phone and new_password are required"})
	}

	err := h.authService.UpdatePassword(req.Phone, req.NewPassword)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.AuthResponse{Message: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.AuthResponse{Message: "Error updating password"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{Message: "Password updated successfully"})
}

func (h *authHandler) UpdatePin(c echo.Context) error {
	var req models.UpdatePinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Invalid request body"})
	}

	if req.Phone == "" || req.NewPin == "" {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "Phone and new_pin are required"})
	}

	if len(req.NewPin) != 6 {
		return c.JSON(http.StatusBadRequest, models.AuthResponse{Message: "PIN must be exactly 6 digits"})
	}

	err := h.authService.UpdatePin(req.Phone, req.NewPin)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.AuthResponse{Message: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.AuthResponse{Message: "Error updating PIN"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{Message: "PIN updated successfully"})
}
