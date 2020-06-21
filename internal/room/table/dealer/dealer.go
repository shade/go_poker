package dealer

import (
	"go_poker/pkg/ringf"
    . "go_poker/internal/interfaces"
    "github.com/chehsunliu/poker"
)


type Dealer struct {
	deck *poker.Deck
	
	flop [3]ICard
	turn ICard
	river ICard
}

func NewDealer() IDealer {
    return &Dealer {
        deck: poker.Deck.NewDeck()
    }
}

func (Dealer *d) DealHands(seats *ringf.RingF) {
	d.deck.Shuffle()

	seats.Do(func(p interface{}) {
		p.(IPlayer).SetHand(d.deck.Draw(2))
	})
}

func (Dealer *d) DealFlop(b Broadcastable) {
	d.flop = d.deck.Draw(3)

	b.Broadcast(msgpb.Packet{
		Event: msgpb.ServerEvent_TABLE_FLOP,
		Payload: 
	})
}

func (Dealer *d) DealTurn(b Broadcastable) {
    return d.deck.Draw(1)[0]
}

func (Dealer *d) DealRiver(b Broadcastable) {
    return d.deck.Draw(1)[0]
}
