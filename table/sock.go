package table

import (
	"fmt"
	"time"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	HB_INTERVAL = 5 * time.Second
	NEW_MSG_DELIMETER = "\n"
	INPUT_BUFFER_SIZE = 1024
	OUTPUT_BUFFER_SIZE = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Sock struct {
	outQ chan []byte
	inQ chan []byte

	connCounter int64
	conns map[int64]*websocket.Conn
}

func NewSock() ISock {
	return &Sock{
		inQ: make(chan []byte, INPUT_BUFFER_SIZE),
		outQ: make(chan []byte, OUTPUT_BUFFER_SIZE),

		connCounter: 0,
		conns: map[int64]*websocket.Conn{},
	}
}

func (s *Sock) read(conn *websocket.Conn, idx int64) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: log error or something
			}

			// TODO: Kill routine for this socket
			conn.Close()
			delete(s.conns, idx)
			break
		}

		s.outQ <- msg
	}
}

func (s *Sock) write(conn *websocket.Conn) {

	for  {
		select {
			case msg := <- s.inQ:
				fullMsg := string(msg) + NEW_MSG_DELIMETER
				for _, conn := range s.conns {
					err := conn.WriteMessage(websocket.TextMessage, []byte(fullMsg))

					if err != nil {
						// TODO: log this error
					}
				}
		}
	}
}

func (s *Sock) AddConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: log errors somehow
		fmt.Println("Error in upgrade")
		fmt.Println(err)
		return
	}

	idx := s.connCounter
	s.conns[idx] = conn
	s.connCounter += 1

	go s.read(conn, idx)
	go s.write(conn)
}

func (s *Sock) Read() []byte {
	return <-s.outQ
}

func (s *Sock) Write(msg []byte) {
	s.inQ <-msg
}
