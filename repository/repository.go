package repository

import (
	"context"

	"github.com/leguminosa/profile-open-portal/entity"
)

//go:generate mockgen -source=repository/repository.go -destination=repository/repository.mock.gen.go -package=repository

type UserRepositoryInterface interface {
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) (int, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	IncrementLoginCount(ctx context.Context, userID int) error
}
