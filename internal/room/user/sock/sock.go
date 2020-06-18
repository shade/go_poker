package sock

import (
	"fmt"
	msgpb "go_poker/internal/proto"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	HB_INTERVAL      = 5 * time.Second
	CHAN_BUFFER_SIZE = 1024

	// Format enums
	FORMAT_JSON  = "json"
	FORMAT_PROTO = "proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SocketConn struct {
	*websocket.Conn
	Id     int64
	IsJSON bool
}

func (s *SocketConn) Write(msg proto.Message) error {
	if s.IsJSON {
		str, err := (&jsonpb.Marshaler{}).MarshalToString(proto.MessageV1(msg))
		if err != nil {
			return err
		}
		return s.WriteMessage(websocket.TextMessage, []byte(str))
	} else {
		out, err := proto.Marshal(msg)
		if err != nil {
			return err
		}
		return s.WriteMessage(websocket.BinaryMessage, out)
	}
}

type Sock struct {
	outQ chan proto.Message
	inQ  chan proto.Message

	connCounter int64
	conns       []*SocketConn
}

func NewSock() *Sock {
	return &Sock{
		inQ:  make(chan proto.Message, CHAN_BUFFER_SIZE),
		outQ: make(chan proto.Message, CHAN_BUFFER_SIZE),

		connCounter: 0,
		conns:       []*SocketConn{},
	}
}

func (s *Sock) read(conn *SocketConn) {
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Invalid websocket reaD!")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: log error or something
			}

			// TODO: Kill routine for this socket
			conn.Close()
			break
		}

		var parseErr error
		msg := &msgpb.Packet{}
		// Parse as JSON or raw proto
		if conn.IsJSON {
			parseErr = jsonpb.UnmarshalString(string(msgBytes), msg)
		} else {
			parseErr = proto.Unmarshal(msgBytes, msg)
		}

		if parseErr != nil {
			fmt.Println("Invalid proto!")
			// TODO: handle /log error
		}

		// Though sock is supposed to be opaque handling ping/pong
		// is an appropriate layer violation
		if msg.Event == msgpb.EventType_PING {
			s.Write(&msgpb.Packet{
				Event: msgpb.EventType_PONG,
			})
			continue
		}

		s.outQ <- msg
	}
}

func (s *Sock) write(conn *SocketConn) {
	for {
		msg := <-s.inQ
		for _, conn := range s.conns {
			err := conn.Write(msg)

			if err != nil {
				// TODO: log this error
			}
		}
	}
}

func (s *Sock) AddConnection(w http.ResponseWriter, r *http.Request) {
	formatArr, ok := r.URL.Query()["format"]
	var format string
	// Default to JSON
	if !ok {
		format = "json"
	} else {
		format = formatArr[0]
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: log errors somehow
		fmt.Println("Error in upgrade")
		fmt.Println(err)
		return
	}

	sc := &SocketConn{
		Conn: conn,
		Id:   s.connCounter,
	}
	s.connCounter += 1

	switch format {
	case FORMAT_JSON:
		sc.IsJSON = true
	case FORMAT_PROTO:
		sc.IsJSON = false
	default:
		// TODO: log errors
		return
	}

	s.conns = append(s.conns, sc)

	go s.read(sc)
	go s.write(sc)
}

func (s *Sock) Read() proto.Message {
	return <-s.outQ
}

func (s *Sock) Write(msg proto.Message) {
	go func() {
		s.inQ <- msg
	}()
}
