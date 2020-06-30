// interfaces is used as a seperate file to prevent cylical dependencies
// between packages that use concrete implementations
package interfaces

import (
	"net/http"

	"github.com/golang/protobuf/proto"
)

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
	// Add an http connection to the sockets listening
	AddConnection(w http.ResponseWriter, r *http.Request)
}

type Sendable interface {
	Send(proto.Message)
}

type Broadcastable interface {
	Broadcast(proto.Message)
}

type IUser interface {
	Sendable
	GetID() string
	RegisterObserver(proto.GeneratedEnum, func(IUser, proto.Message))
	IgnoreUser(proto.GeneratedEnum)
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
	ShowHand(Broadcastable)
	User() IUser

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
	Players() []IPlayer
}
