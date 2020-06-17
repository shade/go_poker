package room

import (
	msgpb "go_poker/internal/proto"
	"go_poker/pkg/pausabletimer"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type User struct {
	ISock
	id string

	balance int
	timer   *pausabletimer.PausableTimer
}

func NewUser(id string) IUser {
	return &User{
		NewSock()
		id:   id,
	}
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) Send(msg proto.Message) {
	m := jsonpb.Marshaler{}
	result, _ := m.MarshalToString(msg)
	u.sock.Write([]byte(result))
}

func (p User) GetSock() ISock {
	return u.sock
}

func (p User) GetBalance(amount int32) int32 {
	return u.balance
}

func (u *User) SubBalance(amount int32) {
	u.balance -= amount
}

func (u *User) Serialize() proto.Message {
	return
}

func (u *User) RegisterObserver(msgType msgpb.EventType, msg proto.Message) {

}
