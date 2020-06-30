package room

import (
	. "gopoker/internal/interfaces"
	msgpb "gopoker/internal/proto"

	"gopoker/internal/room/table"
	"gopoker/internal/room/table/player"

	"time"

	"github.com/golang/protobuf/proto"
)

type Room struct {
	watchers []IUser
	table    ITable
	msgCount uint32
}

func NewRoom(opts *msgpb.TableOptions) *Room {
	r := &Room{
		table:    nil,
		watchers: []IUser{},
		msgCount: 0,
	}

	r.table = table.NewTable(r, opts)

	return r
}

func (r *Room) AddUser(user IUser) {
	r.watchers = append(r.watchers, user)

	user.RegisterObserver(msgpb.ClientEvent_MSG, r.relayChat)
	user.RegisterObserver(msgpb.ClientEvent_SIT_DOWN, r.seatPlayer)
}

func (r *Room) relayChat(user IUser, packet proto.Message) {
	msg := packet.(*msgpb.ClientPacket).GetChat()
	r.msgCount += 1

	r.Broadcast(&msgpb.ServerPacket{
		Event: msgpb.ServerEvent_PLAYER_MSG,
		Payload: &msgpb.ServerPacket_Chat{
			Chat: &msgpb.PlayerMessage_Chat{
				PlayerId:  user.GetID(),
				Data:      msg.Data,
				MessageId: r.msgCount,
				Timestamp: uint32(time.Now().Unix()),
			},
		},
	})
}

func (r *Room) seatPlayer(u IUser, packet proto.Message) {
	msg := packet.(*msgpb.ClientPacket).GetSit()
	seat := msg.GetSeat()
	buyin := msg.GetBuyin()

	err := r.table.SeatPlayer(player.NewPlayer(u, buyin, seat))

	if err != nil {
		u.Send(&msgpb.ServerPacket{
			Event: msgpb.ServerEvent_PLAYER_SIT_REJECT,
			Payload: &msgpb.ServerPacket_SitReject{
				SitReject: &msgpb.PlayerMessage_SitReject{
					Reason: err.Error(),
				},
			},
		})
	} else {
		r.Broadcast(&msgpb.ServerPacket{
			Event: msgpb.ServerEvent_PLAYER_SIT,
			Payload: &msgpb.ServerPacket_PlayerSit{
				PlayerSit: &msgpb.PlayerMessage_Sit{
					PlayerId: u.GetID(),
					SeatNum:  seat,
					Buyin:    buyin,
				},
			},
		})
	}
}

func (r Room) allUsers() []IUser {
	tblUsers := []IUser{}

	for _, player := range r.table.Players() {
		tblUsers = append(tblUsers, player.User())
	}

	return append(r.watchers, tblUsers...)
}

func (r *Room) FindUser(userId string) IUser {
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
