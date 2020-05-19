package table

type Player struct {
	id string
	sock ISock
}

func NewPlayer(id string) IPlayer {
	return &Player{
		id: id,
		sock: NewSock(),
	}
}

func (p *Player) GetID() string {
	return p.id
}

func (p *Player) Send(msg string) {
	p.sock.Write([]byte(msg))
}

func (p *Player) GetSock() ISock {
	return p.sock
}