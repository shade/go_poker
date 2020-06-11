package table

import (
    "fmt"
    "time"

    "poker_backend/messages"
    "github.com/golang/protobuf/jsonpb"
)

type Table struct {
    id string
    players []IPlayer

    minBuy int
    maxSeats int
    bigBlind int

    msgCounter int32
}

func NewTable(id string, minBuy int, maxSeats int) ITable {
    t := &Table{
        id: id,
        minBuy: minBuy,
        maxSeats: maxSeats,
    }

    return t
}

func (t* Table) FindSeat(p IPlayer) int {
    if len(t.players) >= t.maxSeats {
        p.Send(&messages.Packet{
            Event: messages.EventType_TABLE_SIT_ACK,
            Payload: &messages.Packet_SitAck{
                SitAck: &messages.SitAck {
                    TableId: t.id,
                    SatDown: false, 
                    Reason: "Too many people at this table",
                },
            },
        })
    }

    // Bubble insertion
    id := p.GetID()
    seat := 0
    for i, player := range t.players {
        if player.GetID() < id {
            seat = i
            break
        }
    }

    // Shift over the players and add this one in.
    t.players = append(t.players, nil)
    copy(t.players[(seat + 1):], t.players[seat:])
    t.players[seat] = p

    go t.watchPlayer(p)
    p.Send(&messages.Packet{
        Event: messages.EventType_TABLE_SIT_ACK,
        Payload: &messages.Packet_SitAck{
            SitAck: &messages.SitAck {
                TableId: t.id,
                SatDown: true,
                SeatNum: int32(seat), 
            },
        },
    })

    return seat
}

func (t* Table) Stand(p IPlayer) int {
    id := p.GetID()
    seat := -1

    for i, player := range t.players {
        if player.GetID() == id {
            seat = i
            break
        }
    }

    // TODO: REFACTOR! MAKE THIS A FACTORY GROSS.
    t.Broadcast(&messages.Packet{
        Event: messages.EventType_TABLE_STAND_ACK,
        Payload: &messages.Packet_StandAck{
            StandAck: &messages.StandAck {
                TableId: t.id,
                StoodUp: true,
                Balance: 0,
                Reason: "",
            },
        },
    });

    // Remove player
    t.players = append(t.players[:seat], t.players[seat+1:]...)

    return 1
}

func (t *Table) watchPlayer(p IPlayer) {
    for {
        msg := p.GetSock().Read() 
        packet := &messages.Packet{}
        err := jsonpb.UnmarshalString(string(msg), packet)
        if err != nil {
            // TODO: Log this error better
            fmt.Println("Invalid proto receieved")
            continue
        }

        switch packet.Event {
        case messages.EventType_TABLE_STAND:
            t.Stand(p)
        case messages.EventType_CHAT_MSG_SEND:
            t.RelayChat(p, packet.GetPayload().(*messages.Packet_MsgSend))
        default:
            fmt.Println("Invalid event")

        }
    }
}

func (t* Table) Broadcast(packet *messages.Packet) {
    for _, p := range t.players {
        p.Send(packet)
    }
}

func (t* Table) Serialize() *messages.Packet {
    seats := []*messages.PlayerSeat{}

    for i,p := range t.players {
        seats = append(seats, &messages.PlayerSeat{
            Player: p.GetID(),
            Balance: p.GetBalance(),
            SeatNum: int32(i),
        })
    }

    return &messages.Packet{
        Event: messages.EventType_TABLE_STATE,
        Payload: &messages.Packet_JoinState{
            JoinState: &messages.TableState{
                TableId: t.id,
                MinBuy: int32(t.minBuy),
                MaxSeats: int32(t.maxSeats),
                BigBlind: int32(t.bigBlind),
                Seats: seats,
            },
        },
    }
}

func (t* Table) RelayChat(p IPlayer, msg *messages.Packet_MsgSend) {
    data := msg.MsgSend.Data
    player := p.GetID()

    t.Broadcast(&messages.Packet{
        Event: messages.EventType_CHAT_MSG_RECV,
        Payload: &messages.Packet_MsgRecv{
            MsgRecv: &messages.ChatMsgRecv{
                MessageId: t.msgCounter,
                PlayerId: player,
                Data: data,
                Timestamp: int32(time.Now().Unix()),
            },
        },
    })

    t.msgCounter += 1
}