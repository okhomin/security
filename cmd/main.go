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
	log.Println(svc.Signup(context.TODO(), "password", "login3"))
	//log.Println(db.CreateGroup(context.TODO(), group.Group{
	//	Name:  "test",
	//	Read:  true,
	//	Write: true,
	//	Users: []string{
	//		"d5c8a368-b553-449e-aa1e-533682029f96",
	//	},
	//}))
	//log.Println(db.CreateFile(context.TODO(), file.File{
	//	Name:    "Hello",
	//	Content: "testestestes",
	//	Groups: []string{
	//		"de532238-d62b-40e2-9d4c-6d92bf2bf18a",
	//	},
	//	Acls: []string{},
	//}))
	//log.Println(db.CreateAcl(context.TODO(), "d5c8a368-b553-449e-aa1e-533682029f96", true, true))
	log.Println(db.FileGroupPermissions(context.TODO(), "fae54b33-44f0-40b4-b56a-54a80e99ea2f", "d5c8a368-b553-449e-aa1e-533682029f96"))
	log.Println(db.FileGroupPermissions(context.TODO(), "fae54b33-44f0-40b4-b56a-54a80e99ea2f", "d5c8a368-b553-449e-aa1e-533682029f96"))
	//log.Println(db.UpdateFile(context.TODO(), file.File{
	//	ID:      "fae54b33-44f0-40b4-b56a-54a80e99ea2f",
	//	Name:    "Testtest",
	//	Content: "fdsafdsa11111111111",
	//	Groups: []string{
	//		"eff164eb-d222-4df5-9eb5-b72474ae140e",
	//	},
	//	Acls: []string{
	//		"474866a3-d210-4a70-85c5-0e1c6f0eccf3",
	//	},
	//}))
	log.Println(db.FileAclPermissions(context.TODO(), "fae54b33-44f0-40b4-b56a-54a80e99ea2f", "d5c8a368-b553-449e-aa1e-533682029f96"))
}
