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

func main() {
	cfg := config.New()
	db := postgres.New(context.TODO(), cfg.DBURL)
	hasher := bcrypt.New([]byte(cfg.Pepper), cfg.Cost)
	autherSvc := auther.New(hasher, db, db)
	aclerSvc := acler.New(db, db)
	grouperSvc := grouper.New(db, db)
	filerSvc := filer.New(db, db)
	srv := server.New(cfg, autherSvc, aclerSvc, filerSvc, grouperSvc)
	srv.Setup()
	srv.Run()
}
