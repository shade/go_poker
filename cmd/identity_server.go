package main

import (
	"fmt"
	"go_poker/internal/config"
	"go_poker/internal/identity"
	"go_poker/internal/identity/db"
	"go_poker/internal/server"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	var identityDB db.IIDB

	fmt.Println("Connecting to ID Database...")
	switch config.IdentityDB {
	case config.FLAT_FILE:
		identityDB = db.NewFileDB(config.FilePath)
	default:
		panic("Error! Identity DB type not specified.")
	}

	idgen := identity.IDGen(identityDB, config.IDSecret, jwt.SigningMethodHS256)

	fmt.Println("Starting Identity Server...")
	server.Run(config.Address, idgen)
}
