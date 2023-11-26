package main

import (
	"database/sql"
	"os"

	"github.com/leguminosa/profile-open-portal/core"
	_ "github.com/lib/pq"
)

func main() {
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

	// init http server
	core.NewServer(db, privKey, pubKey).Start()
}
