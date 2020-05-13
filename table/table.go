package table

import (
	"poker_backend/interfaces"
)

type Table struct {
	
	players []interfaces.IPlayer

	minBuy int
	maxSeats int
}

func NewTable(minBuy int, maxSeats int) interfaces.ITable {
	return &Table{}
}
