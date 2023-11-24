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

	// init http server
	core.NewServer(db).Start()
}
