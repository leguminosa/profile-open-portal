package tools

import (
	"github.com/labstack/echo/v4"
)

//go:generate mockgen -source=tools/tools.go -destination=tools/tools.mock.gen.go -package=tools

type AuthInterface interface {
	AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type HashInterface interface {
	HashPassword(password string) ([]byte, error)
	ComparePassword(hashedPassword []byte, password string) error
}

type JWTInterface interface {
	Generate(content interface{}) (string, error)
	Validate(tokenString string) (interface{}, error)
}
