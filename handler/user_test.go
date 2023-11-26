package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/module"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserModule := module.NewMockUserModuleInterface(ctrl)

	assert.NotEmpty(t, NewUserHandler(NewUserHandlerOptions{
		UserModule: mockUserModule,
	}))
}

func TestUserHandler_Register(t *testing.T) {
	s := &UserHandler{}
	tests := []struct {
		name    string
		mockCtx *mockEchoContext
		prepare func(m *module.MockUserModuleInterface)
		want    string
		wantErr bool
	}{
		{
			name: "error bind",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					return assert.AnError
				},
			},
			want:    "{\"message\":\"assert.AnError general error for testing\"}\n",
			wantErr: false,
		},
		{
			name: "error register",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.Fullname = "John Doe"
							v.PhoneNumber = "628123456789"
							v.PlainPassword = "Abcde9!"
						}
					}
					return nil
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().Register(mockCtx.Request().Context(), &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "628123456789",
					PlainPassword: "Abcde9!",
				}).Return(entity.RegisterModuleResponse{}, assert.AnError)
			},
			want:    "{\"message\":\"assert.AnError general error for testing\"}\n",
			wantErr: false,
		},
		{
			name: "bad request",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.Fullname = "John Doe"
							v.PhoneNumber = "6212345"
							v.PlainPassword = "Abcde9!"
						}
					}
					return nil
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().Register(mockCtx.Request().Context(), &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "6212345",
					PlainPassword: "Abcde9!",
				}).Return(entity.RegisterModuleResponse{
					Valid:    false,
					Messages: []string{"phone number must be 10-13 digits"},
					User: &entity.User{
						Fullname:      "John Doe",
						PhoneNumber:   "6212345",
						PlainPassword: "Abcde9!",
					},
				}, nil)
			},
			want:    "{\"message\":\"phone number must be 10-13 digits\"}\n",
			wantErr: false,
		},
		{
			name: "success",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.Fullname = "John Doe"
							v.PhoneNumber = "628123456789"
							v.PlainPassword = "Abcde9!"
						}
					}
					return nil
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().Register(mockCtx.Request().Context(), &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "628123456789",
					PlainPassword: "Abcde9!",
				}).Return(entity.RegisterModuleResponse{
					Valid: true,
					User: &entity.User{
						ID:             1,
						Fullname:       "John Doe",
						PhoneNumber:    "628123456789",
						PlainPassword:  "Abcde9!",
						HashedPassword: "hashed Abcde9!",
					},
				}, nil)
			},
			want:    "{\"user_id\":1}\n",
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserModule := module.NewMockUserModuleInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newMockEchoContext(tt.mockCtx)

			if tt.prepare != nil {
				tt.prepare(mockUserModule)
			}
			s.userModule = mockUserModule

			err := s.Register(c)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}
