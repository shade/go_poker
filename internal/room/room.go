package room

import (
	msgpb "go_poker/internal/proto"
	"go_poker/internal/room/user"

	"github.com/golang/protobuf/proto"
	. "go_poker/internal/interfaces"
)

type Room struct {
	watchers []*user.User
	table    *ITable
	msgCount int64
}

func NewRoom(opts table.Options) *Room {
	r = &Room{
		table:    nil,
		watchers: []IPlayer{},
		msgCount: 0,
	}

	r.table = table.NewTable(r, opts)

	return r
}

func (r *Room) AddUser(user *user.User) {
	r.watchers = append(r.watchers, user)

	user.RegisterObserver(msgpb.EventType_CHAT_MSG_SEND, r.relayChat)
	user.RegisterObserver(msgpb.EventType_TABLE_SIT, r.seatPlayer)
}

func (r *Room) relayChat(user *User, packet msgpb.Packet) {
	msg := packet.GetMsgSend()

	r.Broadcast(msgpb.Packet{
		Event: msgpb.EventType_CHAT_MSG_RECV,
		MsgRecv: msgpb.ChatMsgRecv{
			MessageId: r.msgCount++,
			UserId:  user.GetID(),
			Data: msg.Data,
			Timestamp: int32(time.Now().Unix())
		}
	})
}
func (r *Room) seatPlayer(user *User, packet msgpb.Packet) {
	sit := packet.GetSit()

	if r.table.IsValidBuyin(sit.Chips) {
		r.table.Sit(user)

		r.Broadcast(msgpb.Packet{
			Event: msgpb.EventType_TABLE_SIT_ACK,
			SitAck: msgpb.SitAck{
				SatDown: true,
			},
		})
	} else {
		u.Send(msgpb.Packet{
			Event: msgpb.EventType_TABLE_SIT_ACK,
			SitAckt: msgpb.SitAck{
				SatDown: false,
				Reason: "Invalid number of chips",
			}
		})
	}






}

func (r Room) allUsers() []*user.User {
	return append(r.watchers, r.table.GetPlayers()...)
}

func (r *Room) FindUser(string userId) *user.User {
	for _, user := range r.allUsers() {
		if user.GetID() == userId {
			return user
		}
	}

	return nil
}

func (r *Room) Broadcast(msg proto.Message) {
	for _, user := range r.allUsers() {
		user.Send(msg)
	}
}
