package entity

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leguminosa/profile-open-portal/tools"
	"github.com/stretchr/testify/assert"
)

func TestUser_Exist(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want bool
	}{
		{
			name: "user does not exist",
			user: &User{
				ID: 0,
			},
			want: false,
		},
		{
			name: "user exists",
			user: &User{
				ID: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.Exist()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUser_HashPassword(t *testing.T) {
	tests := []struct {
		name               string
		user               *User
		prepare            func(m *tools.MockHashInterface)
		wantHashedPassword string
		wantErr            bool
	}{
		{
			name: "error",
			user: &User{
				PlainPassword: "abcde",
			},
			prepare: func(m *tools.MockHashInterface) {
				m.EXPECT().HashPassword("abcde").Return(nil, assert.AnError)
			},
			wantHashedPassword: "",
			wantErr:            true,
		},
		{
			name: "success",
			user: &User{
				PlainPassword: "abcde",
			},
			prepare: func(m *tools.MockHashInterface) {
				m.EXPECT().HashPassword("abcde").Return([]byte("hashed abcde"), nil)
			},
			wantHashedPassword: "hashed abcde",
			wantErr:            false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHash := tools.NewMockHashInterface(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(mockHash)
			}

			err := tt.user.HashPassword(mockHash)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantHashedPassword, tt.user.HashedPassword)
		})
	}
}
