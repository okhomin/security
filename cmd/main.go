package main

import (
	"context"
	"log"

	"github.com/okhomin/security/internal/models/group"

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
	log.Println(svc.Signup(context.TODO(), "password", "login3"))

	log.Println(db.CreateGroup(context.TODO(), group.Group{
		Name:  "hello",
		Read:  true,
		Write: false,
		Users: []string{"login0", "login1", "test"},
	}))
	log.Println(db.Permissions(context.TODO(), "hello", "12a33cbb-d7b8-47bd-bcb9-8614e1e92089"))
}
