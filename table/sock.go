package table

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	HB_INTERVAL = 5 * time.Second
	NEW_MSG_DELIMETER = '\n'
	INPUT_BUFFER_SIZE = 1024
	OUTPUT_BUFFER_SIZE = 1024
)

type Sock struct {
	outQ chan []byte
	inQ chan []byte

	connCounter int64
	conns map[int64]*websocket.Conn
}

func NewSock() ISock {
	return &Sock{
		inQ: make(chan []byte, INPUT_BUFFER_SIZE)
		outQ: make(chan []byte, OUTPUT_BUFFER_SIZE)

		connCounter: 0,
		conns: []map{}
	}
}

func (s *Sock) read(conn *websocket.Conn, idx int64) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			// TODO: Kill routine for this socket
			conn.Close()
			delete(s.conns, idx)
			break
		}

		outQ <- msg
	}
}

func (s *Sock) write(conn *websocket.Conn) {
	hb := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for  {
		select {
			case msg := <- inQ:
				fullMsg = string(msg)

				for i := 0; i < len(inQ); i++ {
					fullMsg += NEW_MSG_DELIMETER
					fullMsg += string(<-inQ)
				}

				for _, conn := range s.conn {
					conn.WriteMessage(websocket.TextMessage, []byte(fullMsg))
				}

				if err := w.Close(); err != nil {
					// TODO: log error somewhere
					return
				}
			case <-hb.C:
				conn.SetWriteDeadline(time.Now().Add(HB_INTERVAL))
				err := conn.WriteMessage(websocket.PingMessage, nil)

				if err != nil {
					return
				}
		}
	}
}

func (s *Sock) AddConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: log errors somehow
		return
	}

	idx := s.connCounter++
	s.conns[idx] = conn

	go s.read(conn, idx)
}

func (s *Sock) Read() []byte {
	return <-outQ
}

func (s *Sock) Write(msg []byte) {
	inQ <-msg
}
