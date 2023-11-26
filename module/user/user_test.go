package user

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/repository"
	"github.com/leguminosa/profile-open-portal/tools"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repository.NewMockUserRepositoryInterface(ctrl)

	assert.NotEmpty(t, New(NewUserModuleOptions{
		UserRepository: mockUserRepo,
	}))
}

func TestUserModule_Register(t *testing.T) {
	ctx := context.Background()
	m := &UserModule{}
	tests := []struct {
		name        string
		user        *entity.User
		prepareHash func(m *tools.MockHashInterface)
		prepareRepo func(m *repository.MockUserRepositoryInterface)
		want        entity.RegisterModuleResponse
		wantErr     bool
	}{
		{
			name: "request is empty",
			user: &entity.User{
				Fullname:      "",
				PhoneNumber:   "",
				PlainPassword: "",
			},
			want: entity.RegisterModuleResponse{
				Valid: false,
				Messages: []string{
					"phone number must be 10-13 digits",
					"phone number must start with 62",
					"full name must be 3-60 characters",
					"password must be 6-64 characters",
					"password must contain at least 1 uppercase letter, 1 number, and 1 special character",
				},
				User: &entity.User{
					Fullname:      "",
					PhoneNumber:   "",
					PlainPassword: "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid phone number",
			user: &entity.User{
				Fullname:      "John Doe",
				PhoneNumber:   "123456",
				PlainPassword: "Abcde3#",
			},
			want: entity.RegisterModuleResponse{
				Valid: false,
				Messages: []string{
					"phone number must be 10-13 digits",
					"phone number must start with 62",
				},
				User: &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "123456",
					PlainPassword: "Abcde3#",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid full name",
			user: &entity.User{
				Fullname:      "Jo",
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			want: entity.RegisterModuleResponse{
				Valid: false,
				Messages: []string{
					"full name must be 3-60 characters",
				},
				User: &entity.User{
					Fullname:      "Jo",
					PhoneNumber:   "62812345678",
					PlainPassword: "Abcde3#",
				},
			},
			wantErr: false,
		},
		{
			name: "error hash password",
			user: &entity.User{
				Fullname:      "John Doe",
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().HashPassword("Abcde3#").Return(nil, assert.AnError)
			},
			want: entity.RegisterModuleResponse{
				Valid:    true,
				Messages: []string{},
				User: &entity.User{
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "",
					PlainPassword:  "Abcde3#",
				},
			},
			wantErr: true,
		},
		{
			name: "error insert user",
			user: &entity.User{
				Fullname:      "John Doe",
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().HashPassword("Abcde3#").Return([]byte("hashed something"), nil)
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().InsertUser(ctx, &entity.User{
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
					PlainPassword:  "Abcde3#",
				}).Return(0, assert.AnError)
			},
			want: entity.RegisterModuleResponse{
				Valid:    true,
				Messages: []string{},
				User: &entity.User{
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
					PlainPassword:  "Abcde3#",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			user: &entity.User{
				Fullname:      "John Doe",
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().HashPassword("Abcde3#").Return([]byte("hashed something"), nil)
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().InsertUser(ctx, &entity.User{
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
					PlainPassword:  "Abcde3#",
				}).Return(1, nil)
			},
			want: entity.RegisterModuleResponse{
				Valid:    true,
				Messages: []string{},
				User: &entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
					PlainPassword:  "Abcde3#",
				},
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHash := tools.NewMockHashInterface(ctrl)
	mockUserRepo := repository.NewMockUserRepositoryInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareHash != nil {
				tt.prepareHash(mockHash)
			}
			m.hash = mockHash

			if tt.prepareRepo != nil {
				tt.prepareRepo(mockUserRepo)
			}
			m.userRepository = mockUserRepo

			got, err := m.Register(ctx, tt.user)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserModule_Login(t *testing.T) {
	ctx := context.Background()
	m := &UserModule{}
	tests := []struct {
		name        string
		user        *entity.User
		prepareRepo func(m *repository.MockUserRepositoryInterface)
		prepareHash func(m *tools.MockHashInterface)
		prepareJWT  func(m *tools.MockJWTInterface)
		want        entity.LoginModuleResponse
		wantErr     bool
	}{
		{
			name: "error get user",
			user: &entity.User{
				PhoneNumber: "99",
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByPhoneNumber(ctx, "99").Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "user is empty",
			user: &entity.User{
				PhoneNumber: "99",
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByPhoneNumber(ctx, "99").Return(&entity.User{}, nil)
			},
			want: entity.LoginModuleResponse{
				User: &entity.User{},
			},
			wantErr: true,
		},
		{
			name: "hash does not match",
			user: &entity.User{
				PhoneNumber:   "62812345678",
				PlainPassword: "wrong password",
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByPhoneNumber(ctx, "62812345678").Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}, nil)
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().ComparePassword([]byte("hashed something"), "wrong password").Return(assert.AnError)
			},
			want: entity.LoginModuleResponse{
				User: &entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				},
			},
			wantErr: true,
		},
		{
			name: "error generate jwt",
			user: &entity.User{
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByPhoneNumber(ctx, "62812345678").Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}, nil)
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().ComparePassword([]byte("hashed something"), "Abcde3#").Return(nil)
			},
			prepareJWT: func(m *tools.MockJWTInterface) {
				m.EXPECT().Generate(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}).Return("", assert.AnError)
			},
			want: entity.LoginModuleResponse{
				User: &entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				},
				JWT: "",
			},
			wantErr: true,
		},
		{
			name: "error increment login count",
			user: &entity.User{
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByPhoneNumber(ctx, "62812345678").Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}, nil)
				m.EXPECT().IncrementLoginCount(ctx, 1).Return(assert.AnError)
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().ComparePassword([]byte("hashed something"), "Abcde3#").Return(nil)
			},
			prepareJWT: func(m *tools.MockJWTInterface) {
				m.EXPECT().Generate(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}).Return("some jwt token", nil)
			},
			want: entity.LoginModuleResponse{
				User: &entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				},
				JWT: "some jwt token",
			},
			wantErr: true,
		},
		{
			name: "success",
			user: &entity.User{
				PhoneNumber:   "62812345678",
				PlainPassword: "Abcde3#",
			},
			prepareRepo: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByPhoneNumber(ctx, "62812345678").Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}, nil)
				m.EXPECT().IncrementLoginCount(ctx, 1).Return(nil)
			},
			prepareHash: func(m *tools.MockHashInterface) {
				m.EXPECT().ComparePassword([]byte("hashed something"), "Abcde3#").Return(nil)
			},
			prepareJWT: func(m *tools.MockJWTInterface) {
				m.EXPECT().Generate(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}).Return("some jwt token", nil)
			},
			want: entity.LoginModuleResponse{
				User: &entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				},
				JWT: "some jwt token",
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := repository.NewMockUserRepositoryInterface(ctrl)
	mockHash := tools.NewMockHashInterface(ctrl)
	mockJWT := tools.NewMockJWTInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareRepo != nil {
				tt.prepareRepo(mockUserRepo)
			}
			m.userRepository = mockUserRepo

			if tt.prepareHash != nil {
				tt.prepareHash(mockHash)
			}
			m.hash = mockHash

			if tt.prepareJWT != nil {
				tt.prepareJWT(mockJWT)
			}
			m.jwt = mockJWT

			got, err := m.Login(ctx, tt.user)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserModule_GetProfile(t *testing.T) {
	ctx := context.Background()
	m := &UserModule{}
	tests := []struct {
		name    string
		userID  int
		prepare func(m *repository.MockUserRepositoryInterface)
		want    *entity.User
		wantErr bool
	}{
		{
			name:   "error get user",
			userID: 1,
			prepare: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByID(ctx, 1).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "success",
			userID: 1,
			prepare: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByID(ctx, 1).Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
					PlainPassword:  "Abcde3#",
				}, nil)
			},
			want: &entity.User{
				ID:             1,
				Fullname:       "John Doe",
				PhoneNumber:    "62812345678",
				HashedPassword: "hashed something",
				PlainPassword:  "Abcde3#",
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := repository.NewMockUserRepositoryInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(mockUserRepo)
			}
			m.userRepository = mockUserRepo

			got, err := m.GetProfile(ctx, tt.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserModule_UpdateProfile(t *testing.T) {
	ctx := context.Background()
	m := &UserModule{}
	tests := []struct {
		name    string
		user    *entity.User
		prepare func(m *repository.MockUserRepositoryInterface)
		wantErr bool
	}{
		{
			name: "error get user",
			user: &entity.User{},
			prepare: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByID(ctx, 0).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "user not exist",
			user: &entity.User{},
			prepare: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByID(ctx, 0).Return(&entity.User{}, nil)
			},
			wantErr: true,
		},
		{
			name: "error update user",
			user: &entity.User{
				ID:       1,
				Fullname: "John Doe Updated",
			},
			prepare: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByID(ctx, 1).Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}, nil)
				m.EXPECT().UpdateUser(ctx, &entity.User{
					ID:             1,
					Fullname:       "John Doe Updated",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			user: &entity.User{
				ID:          1,
				Fullname:    "John Doe Updated",
				PhoneNumber: "62899123123",
			},
			prepare: func(m *repository.MockUserRepositoryInterface) {
				m.EXPECT().GetUserByID(ctx, 1).Return(&entity.User{
					ID:             1,
					Fullname:       "John Doe",
					PhoneNumber:    "62812345678",
					HashedPassword: "hashed something",
				}, nil)
				m.EXPECT().UpdateUser(ctx, &entity.User{
					ID:             1,
					Fullname:       "John Doe Updated",
					PhoneNumber:    "62899123123",
					HashedPassword: "hashed something",
				}).Return(nil)
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := repository.NewMockUserRepositoryInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(mockUserRepo)
			}
			m.userRepository = mockUserRepo

			err := m.UpdateProfile(ctx, tt.user)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
