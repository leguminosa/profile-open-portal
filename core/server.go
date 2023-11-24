package core

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type ServerInterface interface {
	Start()
}

type Server struct {
	echo *echo.Echo
	app  *App
}

func NewServer(
	db *sql.DB,
) ServerInterface {
	server := &Server{
		echo: echo.New(),
		app:  initApp(db),
	}
	registerMiddleware(server.echo)
	registerHandler(server.echo, server.app)

	return server
}

func (s *Server) Start() {
	s.echo.Logger.Fatal(s.echo.Start(":1323"))
}
