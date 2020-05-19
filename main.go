package main

import (
	"fmt"
	"flag"

	_"poker_backend/table"
	"poker_backend/server"
)

func main() {
	wsport := flag.Int("wsport", 8081, "The port to serve the websocket server")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *wsport)
	server.RunServer(addr)
}
