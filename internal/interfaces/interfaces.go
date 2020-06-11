package table

import (
    "net/http"
    "poker_backend/messages"

    "github.com/golang/protobuf/proto"
)

type IDeck interface {
    Shuffle(seed uint32)
    GetCard(amount int)
}

type ICard interface {
    GetSuit()
    MarshalJSON() ([]byte, error)
}

type IDealer interface {
    GetHand(player IPlayer) [2]ICard

    DealTable() []ICard
    DealFlop() [3]ICard
    DealTurn() ICard
    DealRiver() ICard
}

type ITable interface {
    FindSeat(p IPlayer) int
    Serialize() *messages.Packet
}

type IPlayer interface {
    GetID() string
    Send(proto.Message)
    GetSock() ISock
    GetBalance() int32
}

type ISock interface {
    AddConnection(w http.ResponseWriter, r *http.Request)
    Read() []byte
    Write(msg []byte)
}
