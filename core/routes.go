package core

import "github.com/labstack/echo/v4"

func registerHandler(e *echo.Echo, app *App) {
	e.POST("/register", app.Register)
	e.POST("/login", app.Login)

	v1 := e.Group("/v1", app.auth.AuthenticateMiddleware)
	v1.GET("/profile", app.GetProfile)
	v1.PUT("/profile", app.UpdateProfile)
}
