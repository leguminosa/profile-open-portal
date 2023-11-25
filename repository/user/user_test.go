// This file contains the repository implementation layer.
package user

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer mockDB.Close()

	assert.NotEmpty(t, New(NewRepositoryOptions{
		DB: mockDB,
	}))
}

func TestUserRepository_GetUserByPhoneNumber(t *testing.T) {
	ctx := context.Background()
	r := &UserRepository{}
	tests := []struct {
		name        string
		phoneNumber string
		prepare     func(m sqlmock.Sqlmock)
		want        *entity.User
		wantErr     bool
	}{
		{
			name: "error",
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT.*FROM users WHERE phone_number = \$1`).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT.*FROM users WHERE phone_number = \$1`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"fullname",
						"phone_number",
						"password",
					}).AddRow(
						1,
						"John Doe",
						"628123456789",
						"hashed-password",
					))
			},
			want: &entity.User{
				ID:             1,
				Fullname:       "John Doe",
				PhoneNumber:    "628123456789",
				HashedPassword: "hashed-password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New()
			if err != nil {
				t.Error(err)
			}
			defer mockDB.Close()

			if tt.prepare != nil {
				tt.prepare(mockSQL)
			}
			r.db = mockDB

			got, err := r.GetUserByPhoneNumber(ctx, tt.phoneNumber)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_InsertUser(t *testing.T) {
	ctx := context.Background()
	r := &UserRepository{}
	tests := []struct {
		name    string
		user    *entity.User
		prepare func(m sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "error begin tx",
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "error query row context",
			user: &entity.User{},
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(nil)
				m.ExpectQuery(`INSERT INTO users.*`).WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "error commit",
			user: &entity.User{},
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(nil)
				m.ExpectQuery(`INSERT INTO users.*`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				m.ExpectCommit().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			user: &entity.User{},
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(nil)
				m.ExpectQuery(`INSERT INTO users.*`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				m.ExpectCommit().WillReturnError(nil)
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New()
			if err != nil {
				t.Error(err)
			}
			defer mockDB.Close()

			if tt.prepare != nil {
				tt.prepare(mockSQL)
			}
			r.db = mockDB

			got, err := r.InsertUser(ctx, tt.user)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
