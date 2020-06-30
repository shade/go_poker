package main

import (
	"gopoker/internal/config"
	"gopoker/internal/identity"
	"gopoker/internal/mediator/concierge"
	"gopoker/internal/mediator/custodian"
	"gopoker/internal/server"
	"os"
)

func main() {
	args := os.Args[1:]
	conf := config.ParseFromPath(args[0])

	idgen := identity.NewIDGen(conf.GetDB(), conf.GetTokenSecret())

	custodian := custodian.NewCustodian(conf.GetCache())
	concierge := concierge.NewConcierge(idgen, conf.GetCache())

	go concierge.Start()
	go server.Run(conf.GetAPIAddress(), idgen, custodian)
}
