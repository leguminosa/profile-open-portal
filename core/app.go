package core

import (
	"database/sql"

	"github.com/leguminosa/profile-open-portal/handler"
	moduleUser "github.com/leguminosa/profile-open-portal/module/user"
	repositoryUser "github.com/leguminosa/profile-open-portal/repository/user"
	"github.com/leguminosa/profile-open-portal/tools/crxpto"
	"github.com/leguminosa/profile-open-portal/tools/jwtx"
)

type App struct {
	*handler.UserHandler
}

func initApp(
	db *sql.DB,
	privKey []byte,
	pubKey []byte,
) *App {
	// tools layer
	hashClient := crxpto.NewBcrypt()
	jwtClient := jwtx.NewSigningMethodRS256(jwtx.NewSigningMethodRS256Options{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	})

	// repository layer
	userRepo := repositoryUser.New(repositoryUser.NewRepositoryOptions{
		DB: db,
	})

	// module layer
	userModule := moduleUser.New(moduleUser.NewUserModuleOptions{
		UserRepository: userRepo,
		Hash:           hashClient,
		JWT:            jwtClient,
	})

	// handler layer
	userHandler := handler.NewUserHandler(handler.NewUserHandlerOptions{
		UserModule: userModule,
	})

	return &App{
		UserHandler: userHandler,
	}
}
