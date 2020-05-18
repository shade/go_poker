package table

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	HB_INTERVAL = 5 * time.Second
)

type Sock struct {
	outQ chan []byte
	inQ chan []byte

	connCounter int64
	conns map[int64]*websocket.Conn
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
			case <-hb.C:
				conn.SetWriteDeadline(time.Now().Add(HB_INTERVAL))
				err := conn.WriteMessage(websocket.PingMessage, nil)

				if err != nil {
					return
				}
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
}

func (s *Sock) Write(msg []byte) {

}