package table

import (
	"./interfaces.go"
)

type Table struct {
	
	players []IPlayers

	minBuy int
	maxSeats int
}

func NewTable(minBuy int, maxSeats int) *ITable {

}
