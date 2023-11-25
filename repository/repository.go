package repository

import (
	"context"

	"github.com/leguminosa/profile-open-portal/entity"
)

//go:generate mockgen -source=repository/repository.go -destination=repository/repository.mock.gen.go -package=repository

type UserRepositoryInterface interface {
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) (int, error)
	IncrementLoginCount(ctx context.Context, userID int) error
}
