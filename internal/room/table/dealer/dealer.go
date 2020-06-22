package dealer

import (
	"go_poker/pkg/ringf"
	. "go_poker/internal/interfaces"
	msgpb "go_poker/internal/proto"
    "go_poker/third_party/poker_hands"
)


type Dealer struct {
	deck *poker.Deck
	
	runs int32
	flops [][]ICard
	turns []ICard
	rivers []ICard
}

func NewDealer() *Dealer {
    return &Dealer {
		deck: poker.Deck.NewDeck(),
		runs: 1,
		flops: [][]ICard,
		turns: []ICard,
		rivers: []ICard,
    }
}

func (d* Dealer) DealHands(seats *ringf.RingF) {
	d.deck.Shuffle()
	d.flops = [][]ICard{}
	d.turns = []ICard{}
	d.rivers = []ICard{}

	seats.Do(func(p interface{}) {
		p.(IPlayer).SetHand(d.deck.Draw(2))
	})
}

func (d* Dealer) DealFlop(b Broadcastable) {
	d.flops = make([][]ICard, d.runs)

	for i := 0; i < d.runs; i++ {
		flop := []*msgpb.Card{}
		for _, card := range d.deck.Draw(3){
			d.flops[i] = append(d.flops[i], card)
			flop = append(flop, card.Serialize())
		}

		b.Broadcast(msgpb.ServerPacket{
			Event: msgpb.ServerEvent_TABLE_FLOP,
			Payload: &msgpb.ServerPacket_Flop{
				Flop: &msgpb.CardSet{
					Cards: flop,
				},
			},
		})
	}
}

func (d* Dealer) DealTurn(b Broadcastable) {
	d.turns = make([]ICard, d.runs)

	for i := 0; i < d.runs; i++ {
		card := d.deck.Draw(1)[0]

		d.turns = append(d.turns, card)

		b.Broadcast(msgpb.ServerPacket{
			Event: msgpb.ServerEvent_TABLE_TURN,
			Payload: &msgpb.ServerPacket_Turn{
				Turn: &msgpb.CardSet{
					Cards: []*msgpb.Card{card.Serialize()}
				},
			},
		})
	}
}

func (d* Dealer) DealRiver(b Broadcastable) {
	d.turns = make([]ICard, d.runs)

	for i := 0; i < d.runs; i++ {
		card := d.deck.Draw(1)[0]

		d.turns = append(d.turns, card)

		b.Broadcast(msgpb.ServerPacket{
			Event: msgpb.ServerEvent_TABLE_RIVER,
			Payload: &msgpb.ServerPacket_River{
				River: &msgpb.CardSet{
					Cards: []*msgpb.Card{card.Serialize()}
				},
			},
		})
	}
}
