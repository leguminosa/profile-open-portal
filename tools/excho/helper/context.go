package helper

import (
	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/tools/converter"
)

func UserIDFromContext(c echo.Context) int {
	return converter.ToInt(c.Get("user_id"))
}

func SetUserIDToContext(c echo.Context, userID interface{}) {
	c.Set("user_id", converter.ToInt(userID))
}
