package main

import (
	"flag"
	"fmt"

	"poker_backend/internal/server"
)

func main() {
	wsport := flag.Int("wsport", 8081, "The port to serve the websocket server")
	secret := flag.String("token_secret", "JoeRules123!", "The secret for creating authentication tokens")

	flag.Parse()

	addr := fmt.Sprintf(":%d", *wsport)
	server.RunServer(addr)
}
