package table


type Table struct {
	ITable

	players []Player

	pot int
}

func (t Table) broadcast(Packet) {
	for _, player := range t.players {
		player.SendAction(action)
	}
}

func (t *Table) Sit(p *Player) {
	// Check if player is already sitting
	for _, player := range t.players {
		// Throw bad if player is already sitting
		if player.Id == p.Id {
			t.broadcast()
			return
		}
	}

	// Tell the player they've sat down
	t.players = append(t.players, p)
	t.broadcast(p)
}

func (t *Table) Stand(p *Player) {
	for i,_ := range t.players {
		if player.Id == p.Id {
			t.players = append(t.players[:i], t.players[(i + 1):]])
		}
	}
}

func (t* Table) PlayerCall() {

}

func (t* Table) PlayerRaise(amount int) {
	// Bet must be at least twice the last bet
	if (amount < (t.lastBet * 2)) {
		// Throw error at raise
	}

	t.lastBet = amount

	t.broadcast(raise, amount)
}

func (t* Table) PlayerShove(amount int) {
	
	if (amount > t.lastBet) {
		t.lastBet = amount
	}

	t.broadcast(raise, amount)
}


func (t* Table) PlayerFold() {

}
