package core

import (
	"database/sql"

	"github.com/leguminosa/profile-open-portal/handler"
	moduleUser "github.com/leguminosa/profile-open-portal/module/user"
	"github.com/leguminosa/profile-open-portal/pkg/crxpto"
	repositoryUser "github.com/leguminosa/profile-open-portal/repository/user"
)

type App struct {
	*handler.UserHandler
}

func initApp(db *sql.DB) *App {
	// pkg layer
	hash := crxpto.NewBcrypt()

	// repository layer
	userRepo := repositoryUser.New(repositoryUser.NewRepositoryOptions{
		DB: db,
	})

	// module layer
	userModule := moduleUser.New(moduleUser.NewUserModuleOptions{
		UserRepository: userRepo,
		Hash:           hash,
	})

	// handler layer
	userHandler := handler.NewUserHandler(handler.NewUserHandlerOptions{
		UserModule: userModule,
	})

	return &App{
		UserHandler: userHandler,
	}
}
