package main

import (
	"flag"
	"poker_backend/server"
	"poker_backend/table"
)

func main() {
	wsport := flag.Int("wsport", 8081, "The port to serve the websocket server")
	flag.Parse()


	tbl := table.NewTable()
	server.Start(wsport, tbl)
}
