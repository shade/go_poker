package main

import (
	"go_poker/internal/identity"
	"go_poker/internal/server"
)

func main() {
	filedb := identity.NewFileDB("./hi.csv")
	id := identity.NewIDGen(filedb, "lol")

	server.Run(":8080", id)
}
