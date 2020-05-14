package table

import (
	"poker_backend/msgs"
)

type Table struct {	
	players []IPlayer

	minBuy int
	maxSeats int
}

func NewTable(minBuy int, maxSeats int) ITable {
	return &Table{
		minBuy: minBuy,
		maxSeats: maxSeats,
	}
}

func (t* Table) FindSeat(IPlayer p) int {
	if len(players) >= maxSeats {
		p.Send("No more seats at this table")
	}

	// Bubble insertion
	id := p.GetID()
	seat := 0
	for i, player := range t.players {
		if player.GetID() < id {
			seat = i
			break
		}
	}

	// Shift over the players and add this one in.
	t.players = append(t.players, nil)
	copy(t.players[(seat + 1):], t.players[seat:])
	t.players[seat] = p

	// Broadcast to all players this one is in.
	for i := range t.players[] {
		// TODO: construct seated player message.
	}

	return seat
}

func (t* Table) Stand(IPlayer p) int {
	id := p.GetID()
	seat := -1

	for i, player := range t.players {
		if player.GetID() == id {
			seat = i
			break
		}
	}

	// Remove player
	append(t.players[:seat], t.players[seat+1:]...)


	// TODO: Send and broadcast stand up messages.
	p.Send();
	t.Broadcast();

	return 
}


func (t* Table) Broadcast(msg *msgs.Msg) {

}