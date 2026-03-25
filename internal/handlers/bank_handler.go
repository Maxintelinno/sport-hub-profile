package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
)

type BankHandler struct {
	bankService services.BankService
}

func NewBankHandler(bankService services.BankService) *BankHandler {
	return &BankHandler{bankService: bankService}
}

func (h *BankHandler) GetBankAccounts(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not authenticated"})
	}

	accounts, err := h.bankService.GetOwnerBankAccounts(userID)
	if err != nil {
		if err.Error() == "unauthorized: only owners can access bank accounts" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, models.BankAccountListResponse{
		Status:         "success",
		Message:        "Bank accounts retrieved successfully",
		HasBankAccount: len(accounts) > 0,
		Data:           accounts,
	})
}

func (h *BankHandler) AddBankAccount(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not authenticated"})
	}

	var req models.AddBankAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	account, err := h.bankService.AddOwnerBankAccount(userID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized: only owners can add bank accounts" {
			status = http.StatusForbidden
		} else if err.Error() == "account_name is required" || err.Error() == "either bank account or promptpay info must be provided" {
			status = http.StatusBadRequest
		}
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	resp := models.BankAccountCreateResponse{
		Status:  "success",
		Message: "Bank account created successfully",
	}
	resp.Data.ID = account.ID

	return c.JSON(http.StatusOK, resp)
}
