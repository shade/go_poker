package server

import (
	"go_poker/internal/identity"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run(address string, id *identity.IDGen) {
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/identity", id.IDHandler).Methods("POST")
	s.HandleFunc("/identity/token", id.TokenHandler).Methods("GET")
	s.HandleFunc("/tables", fetchTables).Methods("GET")
	s.HandleFunc("/table", createTables).Methods("POST")

	srv := &http.Server{
		Addr: address,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())

}

func fetchTables(w http.ResponseWriter, r *http.Request) {

}

func createTables(w http.ResponseWriter, r *http.Request) {
	// TODO: Input validation
	// 			Send out creation request
	//			Send back table connection id
}
