package room

import (
	"github.com/golang/protobuf/proto"
)

type Room struct {
	watchers []*interfaces.IPlayer
	table    *table.Table
}

func NewRoom(opts table.Options) *Room {
	r = &Room{
		table:    nil,
		watchers: []*interfaces.IPlayer{},
	}

	r.table = table.NewTable(r, opts)

	return r
}

func (r *Room) AddPlayer(player IPlayer) {
	r.watchers = append(r.watchers, player)
}

func (r *Room) Broadcast(msg proto.Message) {
	allPlayers := append(r.watchers, r.table.GetPlayers()...)

	for _, player := range allPlayers {
		player.Send(msg)
	}
}
