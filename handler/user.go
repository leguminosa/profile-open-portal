package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/module"
	"github.com/leguminosa/profile-open-portal/tools/excho/helper"
)

type UserHandler struct {
	userModule module.UserModuleInterface
}

type NewUserHandlerOptions struct {
	UserModule module.UserModuleInterface
}

func NewUserHandler(opts NewUserHandlerOptions) *UserHandler {
	return &UserHandler{
		userModule: opts.UserModule,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		user = &entity.User{}
	)

	err := c.Bind(user)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	var result entity.RegisterModuleResponse
	result, err = h.userModule.Register(ctx, user)
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}
	if !result.Valid {
		return helper.BadRequest(c, strings.Join(result.Messages, ", "))
	}

	return helper.OK(c, map[string]interface{}{
		"user_id": result.User.ID,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		user = &entity.User{}
	)

	err := c.Bind(&user)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	var result entity.LoginModuleResponse
	result, err = h.userModule.Login(ctx, user)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.OK(c, map[string]interface{}{
		"user_id": result.User.ID,
		"jwt":     result.JWT,
	})
}
