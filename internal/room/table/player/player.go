package player

import (
	"go_poker/internal/room/user"
	msgpb "go_poker/internal/proto"
)

type Player struct {
	*user.User

	balance uint64
	inPot uint64
	state   msgpb.PlayerState
	seat uint64
}

func NewPlayer(u *user.User, balance ) *Player {
	p := &Player{
		u,
		balance,
		msgpb.PlayerState_PENDING,
	}

	p.RegisterObserver(msgpb.EventType_ACTION, p.ProcessAction)

	if table.IsCreator(u) {
		p.RegisterObserver(msgpb.EventType_START_GAME, p.StartGame)
		p.RegisterObserver(msgpb.EventType_END_GAME, p.EndGame)	
	}

	return p
}

func (p *Player) SetSeat(int32 seat) {
	p.seat = seat
}

func (p *Player) MakeBet(uint64 amount) uint64 {
	if amount >= p.balance {
		return p.Shove()
	} else {
		p.balance -= amount
		p.inPot += amount
		p.state = msgpb.PlayerState_RAISED
	}

	return amount
}

func (p *Player) Shove() uint64{
	val := p.balance
	p.inPot = val
	p.balance = 0

	p.state = msgpb.PlayerState_ALL_IN
	return val
}

func (p *Player) IsAllIn() bool {
	return p.state == msgpb.PlayerState_ALL_IN
}
func (p *Player) IsBusted() bool {
	return p.balance == 0
}

func (p *Player) StartGame(msg *msgpb.ActionMsg) {
	if !p.table.IsStarted() {
		p.table.Start()
	}
}

func (p *Player) EndGame(msg *msgpb.ActionMsg) {
	if p.table.IsStarted() {
		p.table.End()
	}
}

func (p *Player) ProcessAction(msg *msgpb.ActionMsg) {
	p.table.timer.Pause()

	switch msg.Type {
		case msgpb.Call
		case msgpb.ActionType_ALL_IN:
			p.table.BubbleAction(action)
			p.Shove()
		case msgpb.ActionType_CALL:

		case msgpb.ActionType_RAISE:
			// Reject too large and small bets
			if p.table.IsValidBet(p.bet) {
				p.MakeBet(p)	
			}
		default:
			// Rejection message

	}

	p.table.timer.Resume()
}


func (p *Player) Serialize() proto.Message {
	return msgpb.Player {
		Name: p.Name,
		Balance: p.balance,
		State: p.state,
		SeatNum: p.seat,
	}
}
