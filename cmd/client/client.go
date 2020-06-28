package main

import (
	"gopoker/internal/config"
	"os"
)

func main() {
	args := os.Args[1:]
	config.ParseFromPath(args[0])
}
