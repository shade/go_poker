package server

import (
	"fmt"
	"net/http"

	"poker_backend/table"
	"poker_backend/messages"
	"github.com/golang/protobuf/jsonpb"
)

var tbl table.ITable

func LobbyRoutine(player table.IPlayer) {
	for {
		msg := player.GetSock().Read() 
		packet := &messages.Packet{}
		err := jsonpb.UnmarshalString(string(msg), packet)
		if err != nil {
			// TODO: Log this error better
			fmt.Println("Invalid proto receieved")
			continue
		}

		if packet.Event == messages.EventType_TABLE_SIT {
			tbl.FindSeat(player)
			break
		} else {
			fmt.Println("Invalid event!")
		}
	}
}

func RunServer(addr string) {
	playerIndex := map[string]table.IPlayer{}
	tbl = table.NewTable("abc table", 1,2)
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.URL.Query()["token"]

		if !ok {
			// TODO: log this error somewhere
			return
		}

		player, exists := playerIndex[token[0]]
		if !exists {
			player = table.NewPlayer(token[0])
			playerIndex[token[0]] = player

			player.GetSock().AddConnection(w, r)
			go LobbyRoutine(player)
		} else {
			player.GetSock().AddConnection(w, r)
		}
	})

	// TODO: log this error too!
	http.ListenAndServe(addr, nil)
}