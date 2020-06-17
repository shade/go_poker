package main

import (
	"flag"
	"fmt"

	"go_poker    /internal/server"
)

func main() {
	wsport := flag.Int("wsport", 8081, "The port to serve the websocket server")
	secret := flag.String("token_secret", "JoeRules123!", "The secret for creating authentication tokens")

	flag.Parse()

	addr := fmt.Sprintf(":%d", *wsport)
	server.RunServer(addr, *secret)
}
