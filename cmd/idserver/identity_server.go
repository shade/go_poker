package main

import (
	"fmt"
	"gopoker/internal/config"
	"gopoker/internal/identity"
	"gopoker/internal/identity/db"
	"gopoker/internal/server"

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
