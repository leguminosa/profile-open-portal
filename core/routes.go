package core

import "github.com/labstack/echo/v4"

func registerHandler(e *echo.Echo, app *App) {
	e.POST("/register", app.Register)
}
