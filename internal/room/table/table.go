package table

import (
	"container/ring"
	msgpb "go_poker/internal/proto"
	"math/rand"
	"sync"
)

type RoundState int

const (
	// UNSTARTED is the state of the game before the start game event is sent.
	UNSTARTED = iota
	PREFLOP
	POSTFLOP
	TURN
	RIVER
)

type Table struct {
	id        string
	pending   []IPlayer
	opts      *msgpb.TableOptions
	room      interfaces.Room
	dealer    IDealer
	actionMux sync.Mutex

	roundState RoundState
	// Pointer on the player that has the lowest seat number.
	players *ring.Ring

	// Pointer on the player that holds the button.
	button *ring.Ring
	// Pointer on the player that currently has the action.
	action *ring.Ring
	// Pointer on the player that is the last to make an aggressive move.
	aggressor *ring.Ring
}

func NewTable() ITable {
	return t
}

func (t *Table) SeatPlayer(p IPlayer, seat int) {
	// Allow player to choose any seat
}

func (t *Table) Start() {
	// If this is set, shuffle the seats
	if t.opts.SeatShuffle {
		t.ShuffleSeats()
	}

	// Pick random dealer
	t.button = t.players.Move(rand.Intn(t.players.Len()))
}

func (t *Table) NewRound(delay int) {
	// Remove all the necessary people
	t.KickPlayers()

	// If there's less than 2 people, end the game.
	if t.players.Len() < 2 {
		t.End()
		return
	}
	t.roundState = PREFLOP

	// Move the dealer button
	t.button = t.button.Next()

	// Person to the left of the dealer starts
	t.action = button.Next()
	t.aggressor = t.action

	// Post blinds
	small, big := t.GetBlinds()
	t.button.Prev().Value.(IPlayer).MakeBet(big)
	t.button.Prev().Prev().Value.(IPlayer).MakeBet(small)

	// Deal and wait.
	t.dealer.DealHands(t.players)
	t.AwaitAction()
}

func (t *Table) NextRound() {
	switch t.state {
	case PREFLOP:
		t.dealer.DealFlop(t.room)
		t.AwaitAction()
		t.state = POSTFLOP
	case POSTFLOP:
		t.dealer.DealTurn(t.room)
		t.AwaitAction()
		t.state = TURN
	case TURN:
		t.dealer.DealRiver(t.room)
		t.AwaitAction()
		t.state = RIVER
	case RIVER:
		t.Showdown()
		t.state = PREFLOP
		t.NewRound()
	}
}

func (t *Table) AwaitAction() {
	p := t.action.Value.(IPlayer)
	// Broadcast the actionable player and their actions.
	p.RegisterObserver(msgpb.EventType_ACTION, t.Action)
	t.Expire(func() {
		p.UnregisterObserver(msgpb.EventType_ACTION)
		// Auto fold the player
		// Next Action
		t.MoveAction(p)
	})
}

func (t *Table) MoveAction(p IPlayer) {
	if t.aggressor == t.action.Next() {
		t.NextRound()
		return
	}

	t.action = t.action.Next()
}

func (t *Table) Action(p IPlayer, action msgpb.ActionType, bet uint64) {
	switch action {
	case msgpb.ActionType_FOLD:
		p.Fold()
	case msgpb.ActionType_CALL:
		if t.state == TableState.BETTING {
			// Silently fail
			return
		}

	case msgpb.ActionType_RAISE:
		if !t.IsValidBet(bet) {
			// Silently fail
			return
		}
		p.Bet(bet)
		p.RegisterBet(bet)
	case msgpb.ActionType_ALL_IN:
		amount := p.Shove()
		p.RegisterBet(amount)
	default:
		// Silently fail
		return
	}

	t.MoveAction()
	t.BroadcastState()
	return
}

func (t *Table) Showdown() {
	t.aggressor.Do(func(p interface{}) {
		player := p.(IPlayer)

		if player.IsInHand() {
			player.ShowCards()
		}
	})
}

func (t *Table) KickBusted() {
	// TODO:
	/*
		// Copy the starting pointer
		tr := t.players

		kicked := []IPlayer{}
		for {
			if t.players[i].IsBusted() {
				kicked = append(kicked, t.players[i])
			}
		}

		// Seconds pass is done in case KickPlayer modifies
		// seating of players.
		for _, p := range kicked {
			t.KickPlayer(p, msgpb.KickReason_BUSTED)
		}*/
}

func (t *Table) GetBlinds() (int64, int64) {
	big := t.opts.BigBlind
	return int64(big / 2), big
}

func (t *Table) ShuffleSeats() {
	rand.Shuffle(len(t.players), func(i, j int) {
		t.players[i], t.players[j] = t.players[j], t.players[i]
	})

	for i, p := range t.players {
		p.SetSeat(i)
	}
}

func (t *Table) Expire(p *Player) {
	// Auto fold the player and continue
	p.Action(FOLD)
	t.MoveAction()
}

func (t *Table) IsValidBet(bet uint64) {

}

func (t *Table) IsValidBuyin(amount int64) bool {
	if amount < t.opts.GetMinBuy() {
		return false
	}

	if amount >= t.opts.GetMaxBuy() {
		return false
	}
}
