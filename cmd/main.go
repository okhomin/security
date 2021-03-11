package main

import (
	"context"
	"log"

	"github.com/okhomin/security/internal/hash/bcrypt"
	"github.com/okhomin/security/internal/service/auth"

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
	svc := auth.New(hasher, db, db)
	log.Println(svc.Signup(context.TODO(), []byte("password"), []byte("login")))
	log.Println(svc.Login(context.TODO(), []byte("password"), []byte("login")))
	log.Println(db)
}
