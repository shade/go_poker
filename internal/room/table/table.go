package table

import (
	"container/ring"
	msgpb "go_poker/internal/proto"
	"math/rand"
)

type RoundState int

const (
	PREFLOP  = iota
	POSTFLOP = iota
	TURN     = iota
	RIVER    = iota
)

type Table struct {
	id      string
	pending []IPlayer
	opts    *msgpb.TableOptions
	room    interfaces.Room

	roundState RoundState
	handCount  int32

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

func (t *Table) Start() {
	// If this is set
	if t.opts.SeatShuffle {
		t.ShuffleSeats()
	}

	// Pick random dealer
	t.buttonIdx = rand.Intn(len(t.players))
	t.StartRound()
}

func (t *Table) StartRound() {
	t.roundState = PREFLOP

	t.KickBusted()
	t.MoveButton()
	t.PostBlinds()
	t.DealCards()
	t.AskAction()
}

func (t *Table) NextRound() {
	switch t.state {
	case PREFLOP:
		t.dealer.DealFlop()
		t.AskAction()
	case POSTFLOP:
		t.dealer.DealTurn()
		t.AskAction()
	case TURN:
		t.dealer.DealRiver()
		t.AskAction()
	case RIVER:
		t.Showdown()
	}
}

func (t *Table) Showdown() {
	t.aggressor.Do(func(p interface{}) {
		player := p.(IPlayer)

		if player.IsInHand() {
			player.ShowCards()
		}
	})
}

func (t *Table) KickPlayer(p IPlayer) {
	t.Broadcast()
}

func (t *Table) KickBusted() {
	kicked := []IPlayer{}
	for i := 0; i < len(t.players); i++ {
		if t.players[i].IsBusted() {
			kicked = append(kicked, t.players[i])
		}
	}

	// Seconds pass is done in case KickPlayer modifies
	// seating of players.
	for _, p := range kicked {
		t.KickPlayer(p, msgpb.KickReason_BUSTED)
	}
}

func (t *Table) MoveButton() {
	t.buttonIdx = (t.buttonIdx + 1) % len(t.players)
}

func (t *Table) PostBlinds() {
	smallIdx := (t.buttonIdx + 1) % len(t.players)
	bigIdx := (t.buttonIdx + 2) % len(t.players)

	small, big := t.GetBlinds()

	t.players[smallIdx].MakeBet(small)
	t.players[bigIdx].MakeBet(big)
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

func (t *Table) IsValidBet(bet uint64) {

}

func (t *Table) MoveAction(p IPlayer) {
	t.actionIdx = (t.actionIdx + 1) % len(t.player)

	if t.actionIdx == t.startIdx {
		t.NextRound()
	} else {
		t.timer.Reset()
	}
}

func (t *Table) IsValidBuyin(amount int64) bool {
	if amount < t.opts.GetMinBuy() {
		return false
	}

	if amount >= t.opts.GetMaxBuy() {
		return false
	}
}
