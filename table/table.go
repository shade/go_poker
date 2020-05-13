package Table

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
		p.send("No more seats at this table");
	}
	// TODO: make this deterministic
	// i.e. same person gets the same position each time
	// the only randomization happens every n hands.
}

func (t* Table) Stand(IPlayer p) int {

}


func (t* Table) Broadcast(msg *msgs.Msg) {

}