package player

type Player struct {
	t Table

	id int
	position bool
	chips int
}


func NewPlayer() *Player {
	p := Player{}
	return &p
}


func (p *Player) GetId() int {
	return p.id
}

func (p *Player) Send(string msg) {
	ws.Send(msg)
}

func (p *Player) AddChips(int chips) {
	p.chips += chips
}

func (p *Player) Emit(msg string) {}

