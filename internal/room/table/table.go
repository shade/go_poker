package table

import (
	"errors"
	. "gopoker/internal/interfaces"
	msgpb "gopoker/internal/proto"
	"gopoker/internal/room/table/dealer"
	"gopoker/pkg/pausabletimer"
	"gopoker/pkg/ringf"
	"time"

	"github.com/golang/protobuf/proto"

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

type TableState int

const (
	CHECKABLE TableState = iota
	CALLABLE
)

type Table struct {
	opts *msgpb.TableOptions

	state TableState

	pending    []IPlayer
	room       IRoom
	dealer     *dealer.Dealer
	playersMux sync.Mutex
	actionMux  sync.Mutex

	lastBet uint64
	timer   *pausabletimer.PausableTimer

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

func NewTable(room IRoom, opts *msgpb.TableOptions) ITable {
	return &Table{
		opts: opts,
		room: room,
	}
}

func (t *Table) Broadcast(msg proto.Message) {
	t.players.Do(func(player interface{}) {
		player.(IPlayer).Send(msg)
	})
}

func (t *Table) Start() {
	// If this is set, shuffle the seats
	if t.opts.SeatShuffle {
		t.ShuffleSeats()
	}

	// Pick random dealer
	t.button = t.players.Move(rand.Intn(t.players.Len()))
}

func (t *Table) kickPlayers() {
	t.players = t.players.Filter(func(player interface{}) bool {
		// Check if we're kicking the button.

		kickMsg := &msgpb.ServerPacket{
			Event:   msgpb.ServerEvent_PLAYER_STAND,
			Payload: nil,
		}

		standMsg := &msgpb.ServerPacket_PlayerStand{
			PlayerStand: &msgpb.PlayerMessage_Stand{
				PlayerId: player.(IPlayer).GetID(),
				SeatNum:  player.(IPlayer).Seat(),
				Balance:  player.(IPlayer).Balance(),
				Reason:   msgpb.PlayerMessage_BANNED,
			},
		}

		if player.(IPlayer).IsBusted() {
			standMsg.PlayerStand.Reason = msgpb.PlayerMessage_BUSTED
			kickMsg.Payload = standMsg
			t.room.Broadcast(kickMsg)
			return false
		} else if player.(IPlayer).IsStanding() {
			standMsg.PlayerStand.Reason = msgpb.PlayerMessage_STOOD_UP
			kickMsg.Payload = standMsg
			t.room.Broadcast(kickMsg)
			return false
		}

		return true
	})
}

func (t *Table) End() {
	t.room.BroadcastStatus("GAME OVER")
}

func (t *Table) NewRound(delay uint32) {
	// TODO: Maybe?
	<-time.NewTimer(pausabletimer.MsToDuration(int64(delay))).C
	// Remove all the necessary people
	t.kickPlayers()

	// If there's less than 2 people, end the game.
	if t.players.Len() < 2 {
		t.End()
		return
	}
	t.roundState = PREFLOP

	// Move the dealer button
	t.button = t.button.Next()

	// Person to the left of the dealer starts
	t.action = t.button.Next()
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
		t.resetPlayers()
		t.state = PREFLOP
		t.NewRound(t.opts.GetRoundDelay())
	}

	t.BroadcastState()
}

func (t *Table) resetPlayers() {
	t.players.Do(func(v interface{}) {
		v.(IPlayer).Reset()
	})
}

func (t *Table) SeatPlayer(p IPlayer) error {
	t.playersMux.Lock()
	defer func() { t.playersMux.Unlock() }()

	lowest := t.players
	seat := p.Seat()

	if !t.isValidBuyin(p.Balance()) {
		return errors.New("Invalid buyin.")
	}

	if !t.isValidSeat(seat) {
		return errors.New("Invalid seat, might be taken.")
	}

	if lowest == nil {
		t.players = ringf.New(1)
		t.players.Value = p
		return nil
	}

	if seat < lowest.Value.(IPlayer).Seat() {
		t.players.Prev().Link(&ringf.RingF{
			Value: p,
		})

		t.players = t.players.Prev()
	} else {
		for {
			if seat < lowest.Value.(IPlayer).Seat() {
				lowest.Prev().Link(&ringf.RingF{
					Value: p,
				})
				break
			}
		}
	}

	return nil
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
	p.WatchPlayer(msgpb.ClientEvent_ACTION, t.Action)
	t.Expire(t.opts.GetActionTimeout(), func() {
		if !t.LockableAction(p) {
			return
		}

		p.IgnorePlayer(msgpb.ClientEvent_ACTION)
		// Auto fold the player
		// Next Action
		t.MoveAction()
		t.actionMux.Unlock()
	})
}

func (t *Table) MoveAction() {

	for {
		if t.aggressor == t.action.Next() {
			t.NextRound()
			return
		}

		player := t.action.Next().Value.(IPlayer)
		if player.IsInHand() && !player.IsAllIn() {
			break
		}
		t.action = t.action.Next()
	}

	t.action = t.action.Next()
}

func (t *Table) validAction(actionType msgpb.ClientActionType) bool {
	if actionType == msgpb.ClientActionType_FOLD {
		return true
	}

	switch t.state {
	case CHECKABLE:
		return ((msgpb.ClientActionType_CHECK == actionType) ||
			(msgpb.ClientActionType_BET == actionType) ||
			(msgpb.ClientActionType_ALL_IN == actionType))
	case CALLABLE:
		return ((msgpb.ClientActionType_CALL == actionType) ||
			(msgpb.ClientActionType_RAISE == actionType) ||
			(msgpb.ClientActionType_ALL_IN == actionType))
	}

	return false
}

func (t *Table) rejectAction(p IPlayer, nonce uint32, reason string) {
	p.Send(&msgpb.ServerPacket{
		Event: msgpb.ServerEvent_ACTION_ACK,
		Payload: &msgpb.ServerPacket_ActionAck{
			ActionAck: &msgpb.ActionAck{
				Ok:    false,
				Error: reason,
				Nonce: nonce,
			},
		},
	})
}

func (t *Table) Action(p IPlayer, payload proto.Message) {
	if !t.LockableAction(p) {
		return
	}

	defer func() {
		t.actionMux.Unlock()
	}()

	actionMsg, ok := payload.(*msgpb.ActionMessage)

	if !ok {
		return
	}

	if !t.validAction(actionMsg.Action) {
		t.rejectAction(p, actionMsg.Nonce, "Invalid Action.")
		return
	}

	switch actionMsg.Action {
	case msgpb.ClientActionType_FOLD:
		p.Fold()
	case msgpb.ClientActionType_CHECK:
		// Just ignore, it's a check
	case msgpb.ClientActionType_CALL:
		p.MakeBet(t.lastBet)
	case msgpb.ClientActionType_BET:
		if p.Balance() < actionMsg.Chips {
			t.rejectAction(p, actionMsg.Nonce, "Invalid number of chips")
			return
		}

		_, bigBlind := t.GetBlinds()
		if actionMsg.Chips < bigBlind {
			t.rejectAction(p, actionMsg.Nonce, "Invalid bet, must be at least big blind")
			return
		}

		p.MakeBet(actionMsg.Chips)
	case msgpb.ClientActionType_RAISE:

	case msgpb.ClientActionType_ALL_IN:
		p.Shove()
	default:
		// Silently fail
		return
	}

	t.MoveAction()
	t.BroadcastState()
	return
}

func (t *Table) BroadcastState() {
	t.room.Broadcast(&msgpb.ServerPacket{})
}

func (t *Table) Showdown() {
	// Get the players with best hand
	players := t.dealer.BestHands(t.aggressor)

	// Chop the pot.
	pot := uint64(0)
	t.players.Do(func(p interface{}) {
		pot += p.(IPlayer).InPot()
	})

	payout := (pot / uint64(len(players)))
	remainder := pot / uint64(len(players))

	// Show their cards to everyone and update chip count
	for _, p := range players {
		p.ShowHand(t.room)
		p.AddChips(payout)
	}

	// Give remainder to the most out of position player
	players[len(players)-1].AddChips(remainder)

}

func (t *Table) GetBlinds() (uint64, uint64) {
	big := t.opts.BigBlind
	return uint64(big / 2), big
}

func (t *Table) isValidBet(bet uint64) bool {
	return bet >= (t.lastBet * 2)
}

func (t *Table) isValidSeat(seat uint32) bool {
	if seat >= t.opts.MaxSeats {
		return false
	}

	taken := t.players.Any(func(player interface{}) bool {
		return player.(IPlayer).Seat() == seat
	})

	if taken {
		return false
	}

	return true
}

func (t *Table) Expire(delay uint32, cb func()) {
	t.timer = pausabletimer.New(int64(delay), cb)

}

func (t *Table) isValidBuyin(amount uint64) bool {
	if amount < t.opts.GetMinBuy() {
		return false
	}

	if amount >= (math.MaxInt64 / uint64(t.opts.MaxSeats)) {
		return false
	}

	return true
}

func (t *Table) ShuffleSeats() {
	// TODO
}

func (t *Table) Players() []IPlayer {
	players := []IPlayer{}

	t.players.Do(func(v interface{}) {
		players = append(players, v.(IPlayer))
	})

	return players
}
