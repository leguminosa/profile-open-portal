package module

import (
	"context"

	"github.com/leguminosa/profile-open-portal/entity"
)

//go:generate mockgen -source=module/module.go -destination=module/module.mock.gen.go -package=module

type UserModuleInterface interface {
	Register(ctx context.Context, user *entity.User) (entity.RegisterModuleResponse, error)
	Login(ctx context.Context, user *entity.User) (entity.LoginModuleResponse, error)
}
