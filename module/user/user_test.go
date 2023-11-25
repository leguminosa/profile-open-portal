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
		req         entity.RegisterModuleRequest
		prepareHash func(m *tools.MockHashInterface)
		prepareRepo func(m *repository.MockUserRepositoryInterface)
		want        entity.RegisterModuleResponse
		wantErr     bool
	}{
		{
			name: "request is empty",
			req: entity.RegisterModuleRequest{
				User: &entity.User{
					Fullname:      "",
					PhoneNumber:   "",
					PlainPassword: "",
				},
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
			req: entity.RegisterModuleRequest{
				User: &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "123456",
					PlainPassword: "Abcde3#",
				},
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
			req: entity.RegisterModuleRequest{
				User: &entity.User{
					Fullname:      "Jo",
					PhoneNumber:   "62812345678",
					PlainPassword: "Abcde3#",
				},
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
			req: entity.RegisterModuleRequest{
				User: &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "62812345678",
					PlainPassword: "Abcde3#",
				},
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
			req: entity.RegisterModuleRequest{
				User: &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "62812345678",
					PlainPassword: "Abcde3#",
				},
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
			req: entity.RegisterModuleRequest{
				User: &entity.User{
					Fullname:      "John Doe",
					PhoneNumber:   "62812345678",
					PlainPassword: "Abcde3#",
				},
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

			got, err := m.Register(ctx, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
