package handler

import (
	"github.com/leguminosa/profile-open-portal/module"
	"github.com/leguminosa/profile-open-portal/tools"
)

type Server struct {
	UserModule module.UserModuleInterface
	Auth       tools.AuthInterface
}

type NewServerOptions struct {
	UserModule module.UserModuleInterface
	Auth       tools.AuthInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		UserModule: opts.UserModule,
		Auth:       opts.Auth,
	}
}
