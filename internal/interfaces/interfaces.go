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
	WatchUser(proto.GeneratedEnum, func(IUser, proto.Message))
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
	SetHand([2]ICard)

	WatchPlayer(proto.GeneratedEnum, func(IPlayer, proto.Message))
	IgnorePlayer(proto.GeneratedEnum)
	ShowHand(Broadcastable)
}

// IRoom is the interface for the room itself
// IRoom is responsible for all the players in the
// room but not the table or game itself.
type IRoom interface {
	Broadcastable
	SeatPlayer(IPlayer, int)
}

type ITable interface {
	Broadcastable
	SeatPlayer(IPlayer)
	IsValidBuyin(int64) bool
}
