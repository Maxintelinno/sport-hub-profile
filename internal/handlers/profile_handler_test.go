package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockProfileService struct{}

func (s *mockProfileService) GetProfile(userID string) (*models.ProfileResponse, error) {
	return &models.ProfileResponse{
		User: models.UserSummary{
			Name: "Owner Lastname",
			Role: "เจ้าของสนาม",
		},
	}, nil
}

func TestGetProfile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "1")

	mockSvc := &mockProfileService{}
	h := NewProfileHandler(mockSvc)

	if assert.NoError(t, h.GetProfile(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp models.ProfileResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Owner Lastname", resp.User.Name)
		assert.Equal(t, "เจ้าของสนาม", resp.User.Role)
	}
}
