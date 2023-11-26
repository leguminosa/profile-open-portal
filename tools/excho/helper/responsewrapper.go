package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func OK(c echo.Context, i interface{}) error {
	return JSON(c, http.StatusOK, i)
}

func BadRequest(c echo.Context, message string) error {
	return JSON(c, http.StatusBadRequest, map[string]interface{}{
		"message": message,
	})
}

func Forbidden(c echo.Context, message string) error {
	return JSON(c, http.StatusForbidden, map[string]interface{}{
		"message": message,
	})
}

func Conflict(c echo.Context, message string) error {
	return JSON(c, http.StatusConflict, map[string]interface{}{
		"message": message,
	})
}

func InternalServerError(c echo.Context, message string) error {
	return JSON(c, http.StatusInternalServerError, map[string]interface{}{
		"message": message,
	})
}

func JSON(c echo.Context, code int, i interface{}) error {
	return c.JSON(code, i)
}
