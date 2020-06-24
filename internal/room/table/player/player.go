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
	hand [2]ICard
}

func NewPlayer(u *user.User, balance, seat) *Player {
	p := &Player{
		u,
		balance,
		msgpb.PlayerState_PENDING,
	}


	return p
}

func (p *Player) WatchPlayer(event proto.GeneratedEnum, cb func(IPlayer, proto.Message)) {
	p.RegisterObserver(event, func(i interface{}, msg proto.Message) {
		cb(i.(IPlayer), msg.())
	})
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

func (p *Player) hasChips(chips uint64) bool {
	return p.balance >= chips
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

func (p *Player) SetHand(hand [2]ICard) {
	p.hand = hand

	p.Send(msgpb.Packet {
		Event: msgpb.EventType_HAND,
		Payload: msgpb.Packet_Hand{
			Hand: &msgpb.CardSet {
				Cards: []msgpb.Card{
					hand[0].Serialize(),
					hand[1].Serialize(),
				}
			}
		}	
	})
}

func (p *Player) ShowHand(b Broadcastable) {
	c
}
