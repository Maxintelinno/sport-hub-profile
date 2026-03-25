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

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) ForgotPassword(phone, newPassword string) error {
	args := m.Called(phone, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) CheckPhone(phone string) (bool, error) {
	args := m.Called(phone)
	return args.Bool(0), args.Error(1)
}

func TestCheckPhone(t *testing.T) {
	e := echo.New()

	t.Run("Phone found", func(t *testing.T) {
		reqBody := models.CheckPhoneRequest{Phone: "0123456789"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/check-phone", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc := new(MockAuthService)
		mockSvc.On("CheckPhone", "0123456789").Return(true, nil)
		h := NewAuthHandler(mockSvc)

		if assert.NoError(t, h.CheckPhone(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp models.CheckPhoneResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.IsFound)
			assert.Equal(t, "Phone number is registered", resp.Message)
		}
	})

	t.Run("Phone not found", func(t *testing.T) {
		reqBody := models.CheckPhoneRequest{Phone: "9999999999"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/check-phone", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc := new(MockAuthService)
		mockSvc.On("CheckPhone", "9999999999").Return(false, nil)
		h := NewAuthHandler(mockSvc)

		if assert.NoError(t, h.CheckPhone(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
			var resp models.CheckPhoneResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.IsFound)
			assert.Equal(t, "ไม่มีเบอร์โทรนี้ในระบบ กรุณาตรวจสอบใหม่อีกครั้ง", resp.Message)
		}
	})
}
