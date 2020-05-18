
type Sock struct {
	outQ chan []byte
	inQ chan []byte

	connCounter int64
	conns map[int64]*websocket.Conn
}

func (s *Sock) read(conn *websocket.Conn) {

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}	
		}
	}
}

func (s *Sock) write() {

}

func (s *Sock) AddConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: log errors somehow
		return
	}

	s.conns[s.connCounter++] = conn
	go s.read(conn)
}

func (s *Sock) Read() []byte {
}

func (s *Sock) Write(msg []byte) {

}