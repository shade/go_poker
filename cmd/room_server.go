package main

import (
	"go_poker/internal/identity"
	"go_poker/internal/server"
)

func main() {
	filedb := identity.NewFileDB("./login.csv")
	id := identity.NewIDGen(filedb, "lol")

	server.Run(":8080", id)
}
