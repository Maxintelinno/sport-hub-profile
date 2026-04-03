package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type OwnerStaffHandler struct {
	staffService services.OwnerStaffService
}

func NewOwnerStaffHandler(staffService services.OwnerStaffService) *OwnerStaffHandler {
	return &OwnerStaffHandler{staffService: staffService}
}

func (h *OwnerStaffHandler) GetStaff(c echo.Context) error {
	// Extract user ID from context (set by JWTMiddleware)
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
		})
	}

	staff, err := h.staffService.GetStaffByOwnerID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, staff)
}

func (h *OwnerStaffHandler) DeactivateStaff(c echo.Context) error {
	// Extract user ID from context (set by JWTMiddleware)
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
		})
	}

	staffUserID := c.Param("id")
	if staffUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Staff user ID is required",
		})
	}

	err := h.staffService.DeactivateStaff(userID, staffUserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "record not found" {
			status = http.StatusNotFound
		}
		return c.JSON(status, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Staff member deactivated successfully",
	})
}
