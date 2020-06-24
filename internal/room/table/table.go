package table

import (
	. "go_poker/internal/interfaces"
	msgpb "go_poker/internal/proto"
	"go_poker/internal/room/table/dealer"
	"go_poker/pkg/ringf"

	"math"
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
	id         string
	pending    []IPlayer
	opts       *msgpb.TableOptions
	room       IRoom
	dealer     *dealer.Dealer
	playersMux sync.Mutex
	actionMux  sync.Mutex

	lastBet uint64

	roundState RoundState
	// Pointer on the player that has the lowest seat number.
	// Any mutations to the players list and seating must occur here.
	players *ringf.RingF

	// Pointer on the player that holds the button.
	button *ringf.RingF
	// Pointer on the player that currently has the action.
	action *ringf.RingF
	// Pointer on the player that is the last to make an aggressive move.
	aggressor *ringf.RingF
}

func NewTable(opts) ITable {
	return Table{}
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

func (t *Table) SeatPlayer(p IPlayer, seat uint32, buyin uint64) {
	t.playersMux.Lock()
	lowest := t.players

	if isValidBuyin(buyin) && isValidSeat(seat) {

	}

	if seat < lowest.Value.(IPlayer).GetSeat() {
		t.players.Prev().Link(ringf.RingF{
			Value: p,
		})

		t.players = t.players.Prev()
	} else {
		for lowest.Next() != t.players {
			if seat < lowest.Value.(IPlayer).GetSeat() {
				lowest.Prev().Link(ringf.RingF{
					Value: p,
				})
				return
			}
		}
	}
	t.playersMux.Unlock()
}

func (t *Table) LockableAction(p IPlayer) bool {
	t.actionMux.Lock()
	if t.action.Value.(IPlayer) == p {
		return true
	} else {
		t.actionMux.Unlock()
		return false
	}
}

func (t *Table) AwaitAction() {
	p := t.action.Value.(IPlayer)
	// Broadcast the actionable player and their actions.
	p.RegisterObserver(msgpb.EventType_ACTION, t.Action)
	t.Expire(func() {
		if !t.LockedAction(p) {
			return
		}

		p.UnregisterObserver(msgpb.EventType_ACTION)
		// Auto fold the player
		// Next Action
		t.MoveAction(p)
		p.actionMux.Unlock()
	})
}

func (t *Table) MoveAction(p IPlayer) {
	if t.aggressor == t.action.Next() {
		t.NextRound()
		return
	}

	t.action = t.action.Next()
}

func (t *Table) Action(p IPlayer, payload proto.Messsage) {
	if !t.LockableAction(p) {
		return
	}

	action, ok := payload.(msgpb.Action)

	if !ok {
		// Log error
	}

	switch action {
	case msgpb.ClientActionType_FOLD:
		p.Fold()
	case msgpb.ClientActionType_CHECK:
		if t.roundState == BETTING {
			t.emitAction(false, "")
			t.playersMux.Unlock()
			return
		}
	case msgpb.ClientActionType_CALL:
		if t.roundState != BETTING {

		}
	case msgpb.ClientActionType_BET:
		if t.roundState != BETTING {

		}
		// check or call
		if action.chips == 0 {

		}

		if p.hasChips(action.chips) {

		}
		if t.isValidBet(action.chips) {

		}

		if t.state == TableState.BETTING {
		}

	case msgpb.ClientActionType_ALLIN:
	default:
		// Silently fail
		return
	}

	t.MoveAction()
	t.BroadcastState()
	t.actionMux.Unlock()
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
	t.playersMux.Lock()

	t.players = t.players.Filter(func(p interface{}) bool {
		if p.(IPlayer).IsBusted() {
			t.KickPlayer(p, msgpb.KickReason_BUSTED)
			return false
		} else {
			return true
		}
	})

	t.playersMux.Unlock()
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

func (t *Table) isValidBet(bet uint64) bool {

}

func (t *Table) isValidSeat(seat uint32) bool {
	if seat >= t.opts.MaxSeats {
		return false
	}

	taken := t.players.Any(func(player interface{}) bool {
		return player.(IPlayer).GetSeat() == seat
	})

	if taken {
		return false
	}

	return false
}

func (t *Table) isValidBuyin(amount uint64) bool {
	if amount < t.opts.GetMinBuy() {
		return false
	}

	if amount >= (math.MaxInt64 / t.opts.MaxSeats) {
		return false
	}

	return true
}
