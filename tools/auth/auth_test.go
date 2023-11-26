package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/tools"
	"github.com/leguminosa/profile-open-portal/tools/excho/helper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := tools.NewMockJWTInterface(ctrl)

	assert.NotEmpty(t, New(NewAuthOptions{
		JWT: mockJWT,
	}))
}

func TestAuth_AuthenticateMiddleware(t *testing.T) {
	a := &Auth{}
	tests := []struct {
		name       string
		token      string
		prepare    func(m *tools.MockJWTInterface)
		wantCode   int
		wantUserID int
	}{
		{
			name:       "missing authorization header",
			token:      "",
			wantCode:   http.StatusForbidden,
			wantUserID: 0,
		},
		{
			name:  "invalid token",
			token: "Bearer some-token",
			prepare: func(m *tools.MockJWTInterface) {
				m.EXPECT().Validate("some-token").Return(nil, assert.AnError)
			},
			wantCode:   http.StatusForbidden,
			wantUserID: 0,
		},
		{
			name:  "malformed data",
			token: "Bearer valid_token",
			prepare: func(m *tools.MockJWTInterface) {
				m.EXPECT().Validate("valid_token").Return("invalid data", nil)
			},
			wantCode:   http.StatusForbidden,
			wantUserID: 0,
		},
		{
			name:  "user id not found",
			token: "Bearer valid_token",
			prepare: func(m *tools.MockJWTInterface) {
				m.EXPECT().Validate("valid_token").Return(map[string]interface{}{
					"some_key": 128,
				}, nil)
			},
			wantCode:   http.StatusForbidden,
			wantUserID: 0,
		},
		{
			name:  "success",
			token: "Bearer valid_token",
			prepare: func(m *tools.MockJWTInterface) {
				m.EXPECT().Validate("valid_token").Return(map[string]interface{}{
					"id": 128,
				}, nil)
			},
			wantCode:   http.StatusOK,
			wantUserID: 128,
		},
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return helper.OK(c, map[string]interface{}{
			"user_id": helper.UserIDFromContext(c),
		})
	}, a.AuthenticateMiddleware)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJWT := tools.NewMockJWTInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(mockJWT)
			}
			a.jwtClient = mockJWT

			mockW := httptest.NewRecorder()
			mockR := httptest.NewRequest("GET", "/", nil)

			mockR.Header.Set("Authorization", tt.token)

			e.ServeHTTP(mockW, mockR)
			assert.Equal(t, tt.wantCode, mockW.Code)

			body, err := io.ReadAll(mockW.Body)
			assert.NoError(t, err)

			type DummyResponse struct {
				UserID  int    `json:"user_id,omitempty"`
				Message string `json:"message,omitempty"`
			}
			var dummyResp DummyResponse
			err = json.Unmarshal(body, &dummyResp)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantUserID, dummyResp.UserID)
		})
	}
}

func TestAuth_getJWTFromHeader(t *testing.T) {
	a := &Auth{}
	tests := []struct {
		name    string
		token   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty authorization header",
			token:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid authorization format",
			token:   "broken_token",
			want:    "",
			wantErr: true,
		},
		{
			name:    "success",
			token:   "Bearer valid_token",
			want:    "valid_token",
			wantErr: false,
		},
	}
	e := echo.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockW := httptest.NewRecorder()
			mockR := httptest.NewRequest("GET", "/", nil)

			mockR.Header.Set("Authorization", tt.token)

			got, err := a.getJWTFromHeader(e.NewContext(mockR, mockW))
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
