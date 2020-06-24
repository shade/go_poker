package server

import (
	"go_poker/internal/identity"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run(address string, id *identity.IDGen) {
	r := Routes("/api/v1", address, id)

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

func Routes(apiPrefix string, address string, id *identity.IDGen) http.Handler {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
	t := r.PathPrefix(apiPrefix).Subrouter()
	t.HandleFunc("/tables", fetchTables).Methods("GET")
	t.HandleFunc("/tables", createTables).Methods("POST")
	t.Use(id.Middleware)

	s := r.PathPrefix(apiPrefix).Subrouter()
	s.HandleFunc("/users", id.IDHandler).Methods("POST")
	s.HandleFunc("/users/{username}/token", id.TokenHandler).Methods("POST")

	return r
}

func fetchTables(w http.ResponseWriter, r *http.Request) {
	//
	//
}

func createTables(w http.ResponseWriter, r *http.Request) {
	// TODO: Input validation
	// 			Send out creation request
	//			Send back table connection id
}
