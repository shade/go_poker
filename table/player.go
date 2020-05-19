package table

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	sock ISock
}

func NewPlayer() IPlayer {
	return &Player{
		sock: NewSock(),
	}
}
