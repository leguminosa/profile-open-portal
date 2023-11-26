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
	h := &UserHandler{}
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
			h.userModule = mockUserModule

			err := h.Register(c)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	h := &UserHandler{}
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
							v.PhoneNumber = "628123456789"
							v.PlainPassword = "Abcde9!"
						}
					}
					return nil
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().Login(mockCtx.Request().Context(), &entity.User{
					PhoneNumber:   "628123456789",
					PlainPassword: "Abcde9!",
				}).Return(entity.LoginModuleResponse{}, assert.AnError)
			},
			want:    "{\"message\":\"assert.AnError general error for testing\"}\n",
			wantErr: false,
		},
		{
			name: "success",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.PhoneNumber = "628123456789"
							v.PlainPassword = "Abcde9!"
						}
					}
					return nil
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().Login(mockCtx.Request().Context(), &entity.User{
					PhoneNumber:   "628123456789",
					PlainPassword: "Abcde9!",
				}).Return(entity.LoginModuleResponse{
					User: &entity.User{
						ID:             1,
						Fullname:       "John Doe",
						PhoneNumber:    "628123456789",
						PlainPassword:  "Abcde9!",
						HashedPassword: "hashed Abcde9!",
					},
					JWT: "some-jwt",
				}, nil)
			},
			want:    "{\"jwt\":\"some-jwt\",\"user_id\":1}\n",
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
			h.userModule = mockUserModule

			err := h.Login(c)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	h := &UserHandler{}
	tests := []struct {
		name    string
		mockCtx *mockEchoContext
		prepare func(m *module.MockUserModuleInterface)
		want    string
		wantErr bool
	}{
		{
			name: "error get profile",
			mockCtx: &mockEchoContext{
				mockGet: func(key string) interface{} {
					return 91
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().GetProfile(mockCtx.Request().Context(), 91).Return(nil, assert.AnError)
			},
			want:    "{\"message\":\"assert.AnError general error for testing\"}\n",
			wantErr: false,
		},
		{
			name: "success",
			mockCtx: &mockEchoContext{
				mockGet: func(key string) interface{} {
					return 15
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().GetProfile(mockCtx.Request().Context(), 15).Return(&entity.User{
					ID:          15,
					Fullname:    "John Doe",
					PhoneNumber: "628123456789",
				}, nil)
			},
			want:    "{\"fullname\":\"John Doe\",\"phone_number\":\"628123456789\"}\n",
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
			h.userModule = mockUserModule

			err := h.GetProfile(c)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	h := &UserHandler{}
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
			name: "error update profile",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.Fullname = "John Doe Updated"
							v.PhoneNumber = "628123456799"
						}
					}
					return nil
				},
				mockGet: func(key string) interface{} {
					return 15
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().UpdateProfile(mockCtx.Request().Context(), &entity.User{
					ID:          15,
					Fullname:    "John Doe Updated",
					PhoneNumber: "628123456799",
				}).Return(entity.UpdateProfileModuleResponse{}, assert.AnError)
			},
			want:    "{\"message\":\"assert.AnError general error for testing\"}\n",
			wantErr: false,
		},
		{
			name: "conflicting phone number",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.Fullname = "John Doe Updated"
							v.PhoneNumber = "628123456799"
						}
					}
					return nil
				},
				mockGet: func(key string) interface{} {
					return 15
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().UpdateProfile(mockCtx.Request().Context(), &entity.User{
					ID:          15,
					Fullname:    "John Doe Updated",
					PhoneNumber: "628123456799",
				}).Return(entity.UpdateProfileModuleResponse{
					Conflict: true,
					Message:  "phone number already exist",
				}, nil)
			},
			want:    "{\"message\":\"phone number already exist\"}\n",
			wantErr: false,
		},
		{
			name: "success",
			mockCtx: &mockEchoContext{
				mockBind: func(i interface{}) error {
					switch v := i.(type) {
					case *entity.User:
						if v != nil {
							v.Fullname = "John Doe Updated"
							v.PhoneNumber = "628123456799"
						}
					}
					return nil
				},
				mockGet: func(key string) interface{} {
					return 15
				},
			},
			prepare: func(m *module.MockUserModuleInterface) {
				m.EXPECT().UpdateProfile(mockCtx.Request().Context(), &entity.User{
					ID:          15,
					Fullname:    "John Doe Updated",
					PhoneNumber: "628123456799",
				}).Return(entity.UpdateProfileModuleResponse{}, nil)
			},
			want:    "{\"user_id\":15}\n",
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
			h.userModule = mockUserModule

			err := h.UpdateProfile(c)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}
