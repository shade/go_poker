package dealer

import (
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

func (Dealer *d) DealHands(seats *ring.Ring) {
	d.deck.Shuffle()

	seats.Do(func(p interface{}) {
		p.(IPlayer).SetHand(d.deck.Draw(2))
	})
}

func (Dealer *d) DealFlop(seats i.IRoom) {
	d.flop = d.deck.Draw(3)

	seats.Do(func(p interface{}) {
	})
}

func (Dealer *d) DealTurn(seats *interfaces.Room) {
	d.turn = 
    return d.deck.Draw(1)[0]
}

func (Dealer *d) DealRiver(seats *interfaces.Room) {
    return d.deck.Draw(1)[0]
}
