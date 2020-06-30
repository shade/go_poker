package dealer

import (
	. "gopoker/internal/interfaces"
	msgpb "gopoker/internal/proto"
	"gopoker/pkg/ringf"
	"gopoker/third_party/pokerdeck"
)

type Dealer struct {
	deck *pokerdeck.Deck

	runs   int
	flops  [][]ICard
	turns  []ICard
	rivers []ICard
}

func NewDealer() *Dealer {
	return &Dealer{
		deck:   pokerdeck.NewDeck(),
		runs:   1,
		flops:  [][]ICard{},
		turns:  []ICard{},
		rivers: []ICard{},
	}
}

func (d *Dealer) DealHands(seats *ringf.RingF) {
	d.deck.Shuffle()
	d.flops = [][]ICard{}
	d.turns = []ICard{}
	d.rivers = []ICard{}

	seats.Do(func(p interface{}) {
		hand := d.deck.Draw(2)
		p.(IPlayer).SetHand(hand[0], hand[1])
	})
}

func (d *Dealer) DealFlop(b Broadcastable) {
	d.flops = make([][]ICard, d.runs)

	for i := 0; i < d.runs; i++ {
		flop := []*msgpb.Card{}
		for _, card := range d.deck.Draw(3) {
			d.flops[i] = append(d.flops[i], card)
			flop = append(flop, card.Serialize().(*msgpb.Card))
		}

		b.Broadcast(&msgpb.ServerPacket{
			Event: msgpb.ServerEvent_TABLE_FLOP,
			Payload: &msgpb.ServerPacket_Flop{
				Flop: &msgpb.CardSet{
					Cards: flop,
				},
			},
		})
	}
}

func (d *Dealer) DealTurn(b Broadcastable) {
	d.turns = make([]ICard, d.runs)

	for i := 0; i < d.runs; i++ {
		card := d.deck.Draw(1)[0]

		d.turns = append(d.turns, card)

		b.Broadcast(&msgpb.ServerPacket{
			Event: msgpb.ServerEvent_TABLE_TURN,
			Payload: &msgpb.ServerPacket_Turn{
				Turn: &msgpb.CardSet{
					Cards: []*msgpb.Card{card.Serialize().(*msgpb.Card)},
				},
			},
		})
	}
}

func (d *Dealer) DealRiver(b Broadcastable) {
	d.turns = make([]ICard, d.runs)

	for i := 0; i < d.runs; i++ {
		card := d.deck.Draw(1)[0]

		d.turns = append(d.turns, card)

		b.Broadcast(&msgpb.ServerPacket{
			Event: msgpb.ServerEvent_TABLE_RIVER,
			Payload: &msgpb.ServerPacket_River{
				River: &msgpb.CardSet{
					Cards: []*msgpb.Card{card.Serialize().(*msgpb.Card)},
				},
			},
		})
	}
}
