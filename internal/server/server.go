package server

import (
	"net/http"
	"poker_backend/internal/identity"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandler = func(w http.ResponseWriter, r *http.Request)

func RunServer(address string, secret string) {
	id := identity.NewIdentityGenerator(secret)

	r := mux.NewRouter()
	r.HandleFunc("/identity", idCreator(id)).Methods("GET")
	r.HandleFunc("/tables", tableFetcher).Methods("GET")
	r.HandleFunc("/table", tableCreator).Methods("POST")

	http.Handle("/", r)
}

func tableFetcher(w http.ResponseWriter, r *http.Request) {

}

func tableCreator(w http.ResponseWriter, r *http.Request) {
	// TODO: Input validation
	// 			Send out creation request
	//			Send back table connection id
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func idCreator(gen identity.IdentityGenerator) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		userId = ip + strconv.FormatInt(time.Now().NanoSeconds())

		token, err := gen.NewToken(userId)

		if err != nil {

		}

	}
}
