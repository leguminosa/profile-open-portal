package main

import (
	"database/sql"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/leguminosa/profile-open-portal/generated"
	"github.com/leguminosa/profile-open-portal/handler"
	moduleUser "github.com/leguminosa/profile-open-portal/module/user"
	repositoryUser "github.com/leguminosa/profile-open-portal/repository/user"
	"github.com/leguminosa/profile-open-portal/tools/auth"
	"github.com/leguminosa/profile-open-portal/tools/crxpto"
	"github.com/leguminosa/profile-open-portal/tools/jwtx"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	server := newServer()
	generated.RegisterHandlers(e, server)

	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	// get db connection string from environment variable
	dbDsn := os.Getenv("DATABASE_URL")

	// initiate database connection
	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		panic(err)
	}

	// get private and public key path from environment variable
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("PUBLIC_KEY_PATH")

	// read private and public key from file
	var privKey, pubKey []byte
	privKey, err = os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	pubKey, err = os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}

	// tools layer
	hashClient := crxpto.NewBcrypt()
	jwtClient := jwtx.NewSigningMethodRS256(jwtx.NewSigningMethodRS256Options{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	})
	authClient := auth.New(auth.NewAuthOptions{
		JWT: jwtClient,
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

	return handler.NewServer(handler.NewServerOptions{
		UserModule: userModule,
		Auth:       authClient,
	})
}
