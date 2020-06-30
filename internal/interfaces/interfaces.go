// interfaces is used as a seperate file to prevent cylical dependencies
// between packages that use concrete implementations
package interfaces

import "github.com/golang/protobuf/proto"

type ObserverCallback func(proto.Message)
type ICard interface {
	Serialize() proto.Message
}

type ISock interface {
	// Write a specific proto message to all the connections for this user
	Write(proto.Message)
	// Registers an observer for a specific type
	RegisterObserver(proto.GeneratedEnum, ObserverCallback)
	// Deregisters all observers for a specific event
	DeregisterObservers(proto.GeneratedEnum)
}

type IUser interface {
	GetID() string
	RegisterObserver(proto.GeneratedEnum, func(IUser, proto.Message))
	IgnoreUser(proto.GeneratedEnum)
	Send(proto.Message)
}

type Broadcastable interface {
	Broadcast(proto.Message)
}

// IPlayer is a subtype of IUser giving all the
// poker playing functionality to an IUser
type IPlayer interface {
	IUser
	SetHand(ICard, ICard)
	IsBusted() bool
	IsStanding() bool
	Balance() uint64
	Seat() uint32
	MakeBet(uint64) uint64
	Fold()
	Shove() uint64
	IsInHand() bool
	ShowHand() proto.Message

	WatchPlayer(proto.GeneratedEnum, func(IPlayer, proto.Message))
	IgnorePlayer(proto.GeneratedEnum)
}

// IRoom is the interface for the room itself
// IRoom is responsible for all the players in the
// room but not the table or game itself.
type IRoom interface {
	Broadcastable
}

type ITable interface {
	Broadcastable
	SeatPlayer(IPlayer) error
}
