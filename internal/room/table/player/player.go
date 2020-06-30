package player

import (
	. "gopoker/internal/interfaces"
	msgpb "gopoker/internal/proto"

	"github.com/golang/protobuf/proto"
)

type Player struct {
	IUser

	balance uint64
	state   msgpb.PlayerState
	inPot   uint64
	seat    uint32
	hand    [2]ICard
}

func NewPlayer(u IUser, balance uint64, seat uint32) *Player {
	p := &Player{
		IUser:   u,
		balance: balance,
		state:   msgpb.PlayerState_PENDING,
		inPot:   0,
		seat:    seat,
		hand:    [2]ICard{},
	}

	return p
}

func (p Player) Balance() uint64 {
	return p.balance
}

func (p Player) Seat() uint32 {
	return p.seat
}

func (p *Player) WatchPlayer(event proto.GeneratedEnum, cb func(IPlayer, proto.Message)) {
	p.RegisterObserver(event, func(_ IUser, msg proto.Message) {
		cb(p, msg)
	})
}

func (p *Player) IgnorePlayer(event proto.GeneratedEnum) {
	p.IgnoreUser(event)
}

func (p *Player) IsInHand() bool {
	return !p.IsStanding() && !p.isFolded()
}
func (p *Player) SetSeat(seat uint32) {
	p.seat = seat
}

func (p *Player) MakeBet(amount uint64) uint64 {
	if amount >= p.balance {
		return p.Shove()
	} else {
		p.balance -= amount
		p.inPot += amount
		p.state = msgpb.PlayerState_RAISED
	}

	return amount
}

func (p *Player) Shove() uint64 {
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
func (p *Player) IsStanding() bool {
	return p.state == msgpb.PlayerState_STOOD_UP
}

func (p *Player) Serialize() proto.Message {
	return &msgpb.Player{
		Name:    p.GetID(),
		Balance: p.balance,
		State:   p.state,
		SeatNum: p.seat,
	}
}

func (p *Player) SetHand(left ICard, right ICard) {
	p.hand = [2]ICard{left, right}

	p.Send(&msgpb.ServerPacket{
		Event: msgpb.ServerEvent_HAND,
		Payload: &msgpb.ServerPacket_Hand{
			Hand: p.getHand(),
		},
	})
}

func (p *Player) ShowHand(b Broadcastable) {
	b.Broadcast(&msgpb.ServerPacket{
		Event: msgpb.ServerEvent_PLAYER_SHOW,
		Payload: &msgpb.ServerPacket_PlayerShow{
			PlayerShow: &msgpb.PlayerMessage_Show{
				Cards: p.getHand(),
			},
		},
	})
}

func (p *Player) getHand() *msgpb.CardSet {
	return &msgpb.CardSet{
		Cards: []*msgpb.Card{
			p.hand[0].Serialize().(*msgpb.Card),
			p.hand[1].Serialize().(*msgpb.Card),
		},
	}
}

func (p *Player) Fold() {
	p.state = msgpb.PlayerState_FOLD
}

func (p *Player) isFolded() bool {
	return p.state == msgpb.PlayerState_FOLD
}

func (p *Player) User() IUser {
	return p.IUser
}
