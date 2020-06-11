package room

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type Player struct {
	id   string
	sock ISock

	balance int
}

func NewPlayer(id string) IPlayer {
	return &Player{
		id:   id,
		sock: NewSock(),
	}
}

func (p *Player) GetID() string {
	return p.id
}

func (p *Player) Send(msg proto.Message) {
	m := jsonpb.Marshaler{}
	result, _ := m.MarshalToString(msg)
	p.sock.Write([]byte(result))
}

func (p *Player) GetSock() ISock {
	return p.sock
}

func (p *Player) GetBalance() int32 {
	return int32(p.balance)
}
