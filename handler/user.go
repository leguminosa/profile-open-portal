package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/module"
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
		ctx     = c.Request().Context()
		request = entity.RegisterAPIRequest{
			User: &entity.User{},
		}
		response entity.RegisterAPIResponse
	)

	err := c.Bind(&request)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	var result entity.RegisterModuleResponse
	result, err = h.userModule.Register(ctx, entity.RegisterModuleRequest{
		User: request.User,
	})
	if err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}
	if !result.Valid {
		response.Message = strings.Join(result.Messages, ", ")
		return c.JSON(http.StatusBadRequest, response)
	}
	response.UserID = result.User.ID

	return c.JSON(http.StatusOK, response)
}
