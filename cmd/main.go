package main

import (
	"context"

	"github.com/okhomin/security/internal/config"
	"github.com/okhomin/security/internal/server"

	"github.com/okhomin/security/internal/service/filer"

	"github.com/okhomin/security/internal/service/grouper"

	"github.com/okhomin/security/internal/service/acler"

	"github.com/okhomin/security/internal/hash/bcrypt"
	"github.com/okhomin/security/internal/service/auther"

	"github.com/okhomin/security/internal/storage/postgres"
)

const ( // TODO: change to config
	cost   = 12
	pepper = "4ed2d5e50bc558927558c0043c5753cf" // must be at least 32 characters
	dbURL  = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
)

func main() {
	db := postgres.New(context.TODO(), dbURL)
	hasher := bcrypt.New([]byte(pepper), cost)
	autherSvc := auther.New(hasher, db, db)
	aclerSvc := acler.New(db, db)
	grouperSvc := grouper.New(db, db)
	cfg := config.Config{
		Port:      "8888",
		JWTKey:    "hello",
		RootLogin: "root",
	}
	filerSvc := filer.New(db, db)
	srv := server.New(cfg, autherSvc, aclerSvc, filerSvc, grouperSvc)
	srv.Setup()
	srv.Run()
}
