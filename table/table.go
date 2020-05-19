package table

import (
	"poker_backend/messages"
)

type Table struct {
	id string
	players []IPlayer

	minBuy int
	maxSeats int
}

func NewTable(id string, minBuy int, maxSeats int) ITable {
	return &Table{
		id: id,
		minBuy: minBuy,
		maxSeats: maxSeats,
	}
}

func (t* Table) FindSeat(p IPlayer) int {
	if len(t.players) >= t.maxSeats {
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

	return seat
}

func (t* Table) Stand(p IPlayer) int {
	id := p.GetID()
	seat := -1

	for i, player := range t.players {
		if player.GetID() == id {
			seat = i
			break
		}
	}

	// Remove player
	t.players = append(t.players[:seat], t.players[seat+1:]...)


	// TODO: Send and broadcast stand up messages.
	p.Send("stood up");
	// TODO: REFACTOR! MAKE THIS A FACTORY GROSS.
	t.Broadcast(&messages.Packet{
		Event: messages.EventType_TABLE_STAND_ACK,
		Payload: &messages.Packet_StandAck{
			StandAck: &messages.StandAck {
				TableId: t.id,
				StoodUp: true,
				Balance: 0,
				Reason: "",
			},
		},
	});

	return 1
}


func (t* Table) Broadcast(packet *messages.Packet) {
	// TODO: implement
}