package table

import "poker_backend/interfaces"

type Table struct {
	ITable

	players []IPlayer

	pot int
	maxSeats int
	minBuyin int
}

func NewTable(seats int, buyin int) *Table {
	return Table{
		maxSeats: seats,
		minBuyin: buyin,
	};
}
