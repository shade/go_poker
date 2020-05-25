package table

import (
	"fmt"
	"poker_backend/messages"
	"github.com/golang/protobuf/jsonpb"
)

type Table struct {
	id string
	players []IPlayer

	minBuy int
	maxSeats int
	bigBlind int
}

func NewTable(id string, minBuy int, maxSeats int) ITable {
	t := &Table{
		id: id,
		minBuy: minBuy,
		maxSeats: maxSeats,
	}

	return t
}

func (t* Table) FindSeat(p IPlayer) int {
	if len(t.players) >= t.maxSeats {
		p.Send(&messages.Packet{
			Event: messages.EventType_TABLE_SIT_ACK,
			Payload: &messages.Packet_SitAck{
				SitAck: &messages.SitAck {
					TableId: t.id,
					SatDown: false, 
					Reason: "Too many people at this table",
				},
			},
		})
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

	go t.watchPlayer(p)
	p.Send(&messages.Packet{
		Event: messages.EventType_TABLE_SIT_ACK,
		Payload: &messages.Packet_SitAck{
			SitAck: &messages.SitAck {
				TableId: t.id,
				SatDown: true,
				SeatNum: int32(seat), 
			},
		},
	})

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

	// Remove player
	t.players = append(t.players[:seat], t.players[seat+1:]...)

	return 1
}

func (t *Table) watchPlayer(p IPlayer) {
	for {
		msg := p.GetSock().Read() 
		packet := &messages.Packet{}
		err := jsonpb.UnmarshalString(string(msg), packet)
		if err != nil {
			// TODO: Log this error better
			fmt.Println("Invalid proto receieved")
			continue
		}

		if packet.Event == messages.EventType_TABLE_STAND {
			t.Stand(p)
			break
		} else {
			fmt.Println("Invalid event!")
		}
	}
}

func (t* Table) Broadcast(packet *messages.Packet) {
	// TODO: implement
}