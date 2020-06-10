package server;

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

fn Start() {
	http.HandleFunc("/subscribe", subscribe)
	http.ListenAndServe(":1234", nil)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Websocket errors:", err)
		return
	}

	go AddConnection(conn);
}

func AddConnection(conn websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}
		
	}
}