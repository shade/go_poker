// interfaces is used as a seperate file to prevent cylical dependencies
// between packages that use concrete implementations
package interfaces

import "github.com/golang/protobuf/proto"

type ICard interface {
	Serialize() proto.Message
}

type ISock interface {
	// Write a specific proto message to all the connections for this user
	Write(proto.Message)
	// Registers an observer for a specific type
	RegisterObserver(proto.GeneratedEnum, func(interface{}, proto.Message))
	// Deregisters observer for a specific event
	DeregisterObservers(proto.GeneratedEnum)
}

type IUser interface {
	ISock
	WatchUser(proto.GeneratedEnum, func(IUser, proto.Message))
}

// IPlayer is a subtype of IUser giving all the
// poker playing functionality to an IUser
type IPlayer interface {
	IUser
	SetHand([2]ICard)

	WatchPlayer(proto.GeneratedEnum, func(IPlayer, proto.Message))
}

// IRoom is the interface for the room itself
// IRoom is responsible for all the players in the
// room but not the table or game itself.
type IRoom interface {
	Broadcast(proto.Message)
	SeatPlayer(IPlayer, int)
}

type ITable interface {
	SeatPlayer(IPlayer)
}
