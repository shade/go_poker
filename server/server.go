package server

import (
	"net/http"

	"poker_backend/table"
)

var playerIndex map[string]table.IPlayer

func RunServer(addr string) {
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
			
		} 

		player.GetSock().AddConnection(w, r)
	})

	// TODO: log this error too!
	http.ListenAndServe(addr, nil)
}