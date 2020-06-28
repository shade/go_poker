package main

import (
	"go_poker/internal/config"
	"os"
)

func main() {
	args := os.Args[1:]
	config.ParseFromPath(args[0])
}
