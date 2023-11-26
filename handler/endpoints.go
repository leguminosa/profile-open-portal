package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/generated"
	"github.com/leguminosa/profile-open-portal/tools/excho/helper"
)

func (s *Server) PostLogin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = &generated.LoginRequest{}
	)

	err := c.Bind(req)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	var result entity.LoginModuleResponse
	result, err = s.UserModule.Login(ctx, &entity.User{
		PhoneNumber:   req.PhoneNumber,
		PlainPassword: req.Password,
	})
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.OK(c, generated.LoginResponse{
		Jwt:    result.JWT,
		UserId: int64(result.User.ID),
	})
}

func (s *Server) PostRegister(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = &generated.RegisterRequest{}
	)

	err := c.Bind(req)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	var result entity.RegisterModuleResponse
	result, err = s.UserModule.Register(ctx, &entity.User{
		Fullname:      req.Fullname,
		PhoneNumber:   req.PhoneNumber,
		PlainPassword: req.Password,
	})
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}
	if !result.Valid {
		return helper.BadRequest(c, strings.Join(result.Messages, ", "))
	}

	return helper.OK(c, generated.RegisterResponse{
		UserId: int64(result.User.ID),
	})
}

func (s *Server) GetV1Profile(c echo.Context) error {
	if err := s.Auth.Authenticate(c); err != nil {
		return helper.Forbidden(c, err.Error())
	}

	var (
		ctx    = c.Request().Context()
		userID = helper.UserIDFromContext(c)
	)

	result, err := s.UserModule.GetProfile(ctx, userID)
	if err != nil {
		return helper.Forbidden(c, err.Error())
	}

	return helper.OK(c, generated.GetProfileResponse{
		Fullname:    result.Fullname,
		PhoneNumber: result.PhoneNumber,
	})
}

func (s *Server) PutV1Profile(c echo.Context) error {
	if err := s.Auth.Authenticate(c); err != nil {
		return helper.Forbidden(c, err.Error())
	}

	var (
		ctx    = c.Request().Context()
		req    = &generated.UpdateProfileRequest{}
		userID = helper.UserIDFromContext(c)
	)

	err := c.Bind(req)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	var result entity.UpdateProfileModuleResponse
	result, err = s.UserModule.UpdateProfile(ctx, &entity.User{
		ID:          userID,
		Fullname:    req.Fullname,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		return helper.Forbidden(c, err.Error())
	}
	if result.Conflict {
		return helper.Conflict(c, result.Message)
	}

	return helper.OK(c, generated.UpdateProfileResponse{
		UserId: int64(userID),
	})
}
