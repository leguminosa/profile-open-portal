package auth

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/tools"
	"github.com/leguminosa/profile-open-portal/tools/excho/helper"
)

type Auth struct {
	jwtClient tools.JWTInterface
}

type NewAuthOptions struct {
	JWT tools.JWTInterface
}

func New(opts NewAuthOptions) *Auth {
	return &Auth{
		jwtClient: opts.JWT,
	}
}

var (
	// ErrNotAuthenticated obscures the error message
	// to avoid brute force attack on authentication process
	ErrNotAuthenticated = errors.New("not authenticated")
)

func (a *Auth) AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := a.Authenticate(c)
		if err != nil {
			return helper.Forbidden(c, ErrNotAuthenticated.Error())
		}

		return next(c)
	}
}

func (a *Auth) Authenticate(c echo.Context) error {
	jwtToken, err := a.getJWTFromHeader(c)
	if err != nil {
		return ErrNotAuthenticated
	}

	var dat interface{}
	dat, err = a.jwtClient.Validate(jwtToken)
	if err != nil {
		return ErrNotAuthenticated
	}

	user, ok := dat.(map[string]interface{})
	if !ok {
		return ErrNotAuthenticated
	}

	var userID interface{}
	userID, ok = user["id"]
	if !ok {
		return ErrNotAuthenticated
	}

	helper.SetUserIDToContext(c, userID)
	return nil
}

func (a *Auth) getJWTFromHeader(c echo.Context) (string, error) {
	authorizationHeader := c.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("missing authorization header")
	}

	authorizationHeaderPart := strings.Split(authorizationHeader, "Bearer ")
	if len(authorizationHeaderPart) != 2 {
		return "", errors.New("invalid authorization header")
	}
	bearerToken := authorizationHeaderPart[1]

	return bearerToken, nil
}
