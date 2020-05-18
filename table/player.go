package table

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	sock ISock
}

func NewPlayer(w http.ResponseWriter, r *http.Reques) IPlayer {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODD: attempt to report and log error
		return
	}

	p := &Player{
		conn,
	}

	go p.read()

	return p
}
