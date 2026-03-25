package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBankService struct {
	mock.Mock
}

func (m *MockBankService) GetOwnerBankAccounts(userID string) ([]models.BankAccountResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BankAccountResponse), args.Error(1)
}

func (m *MockBankService) AddOwnerBankAccount(userID string, req models.AddBankAccountRequest) (*models.OwnerBankAccount, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OwnerBankAccount), args.Error(1)
}

func (m *MockBankService) UpdateOwnerBankAccount(userID string, accountID string, req models.UpdateBankAccountRequest) error {
	args := m.Called(userID, accountID, req)
	return args.Error(0)
}

func (m *MockBankService) SetDefaultBankAccount(userID string, accountID string) error {
	args := m.Called(userID, accountID)
	return args.Error(0)
}

func (m *MockBankService) DeleteOwnerBankAccount(userID string, accountID string) error {
	args := m.Called(userID, accountID)
	return args.Error(0)
}

func TestGetBankAccounts(t *testing.T) {
	e := echo.New()
	t.Run("Successful retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/owner/bank-accounts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "owner-123")

		mockSvc := new(MockBankService)
		expectedData := []models.BankAccountResponse{
			{
				ID:                  "acc-1",
				BankCode:            "KBANK",
				AccountNumberMasked: "xxx-x-12345-x",
			},
		}
		mockSvc.On("GetOwnerBankAccounts", "owner-123").Return(expectedData, nil)

		h := NewBankHandler(mockSvc)
		if assert.NoError(t, h.GetBankAccounts(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp models.BankAccountListResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, "success", resp.Status)
			assert.True(t, resp.HasBankAccount)
			assert.Equal(t, 1, len(resp.Data))
			assert.Equal(t, "xxx-x-12345-x", resp.Data[0].AccountNumberMasked)
		}
	})
}

func TestAddBankAccount(t *testing.T) {
	e := echo.New()
	t.Run("Successful addition", func(t *testing.T) {
		reqBody := models.AddBankAccountRequest{
			BankCode:    "KBANK",
			AccountName: "Somsak Jaidee",
			AccountNumber: "1234567890",
			IsDefault:   true,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/owner/bank-accounts", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "owner-123")

		mockSvc := new(MockBankService)
		mockSvc.On("AddOwnerBankAccount", "owner-123", reqBody).Return(&models.OwnerBankAccount{ID: "new-acc-id"}, nil)

		h := NewBankHandler(mockSvc)
		if assert.NoError(t, h.AddBankAccount(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp models.BankAccountCreateResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, "success", resp.Status)
			assert.Equal(t, "new-acc-id", resp.Data.ID)
		}
	})
}

func TestUpdateBankAccount(t *testing.T) {
	e := echo.New()
	t.Run("Successful update", func(t *testing.T) {
		mockSvc := new(MockBankService)
		h := NewBankHandler(mockSvc)

		accountID := "acc-1"
		userID := "owner-123"
		reqBody := models.UpdateBankAccountRequest{BankCode: "SCB", AccountName: "New Name"}
		body, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest(http.MethodPut, "/v1/owner/bank-accounts/"+accountID, bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)
		c.Set("user_id", userID)

		mockSvc.On("UpdateOwnerBankAccount", userID, accountID, reqBody).Return(nil)

		if assert.NoError(t, h.UpdateBankAccount(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "success")
		}
	})
}

func TestSetDefaultBankAccount(t *testing.T) {
	e := echo.New()
	t.Run("Successful set default", func(t *testing.T) {
		mockSvc := new(MockBankService)
		h := NewBankHandler(mockSvc)

		accountID := "acc-1"
		userID := "owner-123"
		
		req := httptest.NewRequest(http.MethodPost, "/v1/owner/bank-accounts/"+accountID+"/set-default", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)
		c.Set("user_id", userID)

		mockSvc.On("SetDefaultBankAccount", userID, accountID).Return(nil)

		if assert.NoError(t, h.SetDefaultBankAccount(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "success")
		}
	})
}

func TestDeleteBankAccount(t *testing.T) {
	e := echo.New()
	t.Run("Successful deletion", func(t *testing.T) {
		mockSvc := new(MockBankService)
		h := NewBankHandler(mockSvc)

		accountID := "acc-1"
		userID := "owner-123"
		
		req := httptest.NewRequest(http.MethodDelete, "/v1/owner/bank-accounts/"+accountID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)
		c.Set("user_id", userID)

		mockSvc.On("DeleteOwnerBankAccount", userID, accountID).Return(nil)

		if assert.NoError(t, h.DeleteBankAccount(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "success")
		}
	})
}
