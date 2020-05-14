package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	"poker_backend/table"
)

func main() {
	wsport := flag.Int("wsport", 8081, "The port to serve the websocket server")
	flag.Parse()

	tbl := table.NewTable(1,1)
}
