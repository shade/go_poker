package table

import (
	msgpb "go_poker/internal/proto"
	"math/rand"
)

type Table struct {
	id             string
	pendingPlayers []IPlayer
	players        []IPlayer
	msgCounter     int32
	opts           *msgpb.TableOptions
	room           interfaces.Room

	handCount int32
	buttonIdx int32
	actionIdx int32
	startIdx  int32
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
	// Remove all empty balance players
	t.KickBusted()
	// Move dealer chip once
	t.MoveButton()
	// Post little and big blind
	t.PostBlinds()
	// Ask for action on the first player and wait
	t.AskAction()
}

func (t *Table) KickBusted() {
	for i := 0; i < len(t.players); i++ {
		if t.players[i].IsBusted() {
			// Broadcast the kick!
			t.players[i].Kick()
			// Remove from the slice
			t.players = append(t.players[:i], t.players[i+1:]...)
			i -= 1
		}
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
