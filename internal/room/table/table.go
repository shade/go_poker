package table

import (
	"errors"
	. "gopoker/internal/interfaces"
	msgpb "gopoker/internal/proto"
	"gopoker/internal/room/table/dealer"
	"gopoker/pkg/pausabletimer"
	"gopoker/pkg/ringf"

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

func (t *Table) end() {
	// TODO: Maybe?
}

func (t *Table) NewRound(delay uint32) {
	// Remove all the necessary people
	t.kickPlayers()

	// If there's less than 2 people, end the game.
	if t.players.Len() < 2 {
		t.end()
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
		t.state = PREFLOP
		t.NewRound(t.opts.GetRoundDelay())
	}
}

func (t *Table) SeatPlayer(p IPlayer) error {
	t.playersMux.Lock()
	defer func() { t.playersMux.Unlock() }()

	lowest := t.players
	seat := p.Seat()

	if t.isValidBuyin(p.Balance()) {
		return errors.New("Invalid buyin.")
	}

	if t.isValidSeat(seat) {
		return errors.New("Invalid seat, might be taken.")
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

		if t.action.Next().Value.(IPlayer).IsInHand() {
			break
		}
		t.action = t.action.Next()
	}

	t.action = t.action.Next()
}

func (t *Table) validAction(actionType msgpb.ClientActionType) bool {
	switch actionType {

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

	actionMsg, ok := payload.(*msgpb.ActionMessage)

	if !ok {
		t.actionMux.Unlock()
		return
	}

	if !t.validAction(actionMsg.Action) {
		t.actionMux.Unlock()
		t.rejectAction(p, actionMsg.Nonce, "Invalid Action.")
		return
	}

	switch actionMsg.Action {
	case msgpb.ClientActionType_FOLD:
		p.Fold()
	case msgpb.ClientActionType_CHECK:
	case msgpb.ClientActionType_CALL:
		p.MakeBet(t.lastBet)
	case msgpb.ClientActionType_BET:
		if p.Balance() < actionMsg.Chips {
			t.rejectAction(p, actionMsg.Nonce, "Invalid number of chips")
			t.actionMux.Unlock()
			return
		}

		_, bigBlind := t.GetBlinds()
		if actionMsg.Chips < bigBlind {
			t.actionMux.Unlock()
			return
		}
	case msgpb.ClientActionType_RAISE:

	case msgpb.ClientActionType_ALL_IN:
		p.Shove()
	default:
		// Silently fail
		return
	}

	t.MoveAction()
	t.BroadcastState()
	t.actionMux.Unlock()
	return
}

func (t *Table) BroadcastState() {
	t.room.Broadcast(&msgpb.ServerPacket{})
}

func (t *Table) Showdown() {
	t.aggressor.Do(func(p interface{}) {
		player := p.(IPlayer)

		if player.IsInHand() {
			player.ShowCards()
		}
	})
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

	return false
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
