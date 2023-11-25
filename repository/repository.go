// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/leguminosa/profile-open-portal/entity"
)

//go:generate mockgen -source=repository/repository.go -destination=repository/repository.mock.gen.go -package=repository

type UserRepositoryInterface interface {
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) (int, error)
}
