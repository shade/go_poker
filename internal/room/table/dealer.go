package dealer

import (
    "go_poker/interfaces"
    "github.com/chehsunliu/poker"
)


type Dealer struct {
    deck *poker.Deck
    hands map[int][2]poker.Card
}

func NewDealer() IDealer {
    return &Dealer {
        deck: poker.Deck.NewDeck(),
        hands: nil
    }
}

func (Dealer *d) DealTable(seats: []int) {
    d.deck.Shuffle()
    d.hands = make(map[int][2]poker.Card)
    for _, idx := range 0..players {
        d.hands[idx] = d.deck.Draw(2)
    }
}

func (Dealer *d) DealFlop() [3]ICard {
    return d.deck.Draw(3)poker.Card
}

func (Dealer *d) DealTurn() ICard {
    return d.deck.Draw(1)[0]
}

func (Dealer *d) DealRiver() ICard {
    return d.deck.Draw(1)[0]
}
func (Dealer *d) GetHand(seatIdx: int) ([]poker.Card, error) {
    hand, ok := d.hands[seatIdx]

    if !ok {
        return (nil, errors.New("No hand dealt, for this seat."))
    }

    return (hand, nil)
}
