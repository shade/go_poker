package server

import (
	"net/http"

	"poker_backend/table"
)

var playerIndex = map[string]*IPlayer

func RunServer(ITable* table) {
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.URL.Query()["token"]

		if !ok {
			// TODO: log this error somewhere
			return
		}

		player, exists := playerIndex[token]
		if !exists {
			// Create the player
			playerIndex[token] = table.NewPlayer(w,r)
		} else {
			// Add this connection to the connectionlist
			player.GetSock().AddConnection(w, r)
		}
	})
}