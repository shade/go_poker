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
	flops  [][]pokerdeck.Card
	turns  []pokerdeck.Card
	rivers []pokerdeck.Card
}

func NewDealer() *Dealer {
	return &Dealer{
		deck:   pokerdeck.NewDeck(),
		runs:   1,
		flops:  [][]pokerdeck.Card{},
		turns:  []pokerdeck.Card{},
		rivers: []pokerdeck.Card{},
	}
}

func (d *Dealer) DealHands(seats *ringf.RingF) {
	d.deck.Shuffle()
	d.flops = [][]pokerdeck.Card{}
	d.turns = []pokerdeck.Card{}
	d.rivers = []pokerdeck.Card{}

	seats.Do(func(p interface{}) {
		hand := d.deck.Draw(2)
		p.(IPlayer).SetHand(hand[0], hand[1])
	})
}

func (d *Dealer) DealFlop(b Broadcastable) {
	d.flops = make([][]pokerdeck.Card, d.runs)

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
	d.turns = make([]pokerdeck.Card, d.runs)

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
	d.turns = make([]pokerdeck.Card, d.runs)

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

func (d *Dealer) BestHands(seats *ringf.RingF) []IPlayer {
	best := []IPlayer{}
	bestHand := int32(0)

	seats.Do(func(p interface{}) {
		left, right := p.(IPlayer).Hand()

		if !p.(IPlayer).IsInHand() {
			return
		}

		handVal := pokerdeck.Evaluate(append(d.flops[0], d.turns[0], d.rivers[0], left.(pokerdeck.Card), right.(pokerdeck.Card)))

		if bestHand < handVal {
			bestHand = handVal
			best = []IPlayer{p.(IPlayer)}
		} else if bestHand == handVal {
			best = append(best, p.(IPlayer))
		}
	})

	return best
}
