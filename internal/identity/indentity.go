package identity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// 1 day in nanoseconds
const TOKEN_DURATION = 86486400000000000
const AUTH_HEADER = "Authorization"
const BEARER_SCHEMA = "Bearer "

type IDGen struct {
	db     IIDB
	secret string
	algo   jwt.SigningMethod
}

func NewIDGen(db IIDB, secret string) *IDGen {
	return &IDGen{
		db:     db,
		secret: secret,
		algo:   jwt.SigningMethodHS256,
	}
}

func (i *IDGen) IsValidPassword(hash string, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func (i *IDGen) FetchRecord(username string, password string) (*Record, error) {
	record, err := i.db.Get(DBKey(username))

	if err != nil {
		return nil, errors.New("User doesn't exist")
	}

	if !i.IsValidPassword(record.PasswordHash, password) {
		return nil, errors.New("Invalid password")
	}

	return record, nil
}

func (i *IDGen) ValidateRecord(r *Record) error {
	// Ensure no user with same username in the db
	if record, _ := i.db.Get(DBKey(r.Username)); record != nil {
		return errors.New("User already exists in DB")
	}

	return nil
}

func (i *IDGen) ParseToken(t string) (bool, map[string]interface{}) {
	claims := jwt.MapClaims{}
	fmt.Println(t)
	token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(i.secret), nil
	})

	if err != nil || !token.Valid {
		return false, nil
	}

	return true, claims
}

func (i *IDGen) CreateToken(r *Record) string {
	token, err := jwt.NewWithClaims(i.algo, jwt.MapClaims{
		"username": r.Username,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	}).SignedString([]byte(i.secret))

	if err != nil {
		log.Fatal(err)
	}

	return token
}

type errorMsg struct {
	Error string `json:"error"`
}

type tokenMsg struct {
	Token string `json:"token"`
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	dataStr, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dataStr)
}

func (i *IDGen) IDHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSON(w, errorMsg{Error: "Could not parse form data"})
		return
	}

	name := r.FormValue("name")
	user := r.FormValue("username")
	pass := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), 10)

	record := &Record{name, user, pass, base64.StdEncoding.EncodeToString(hash)}
	if err := i.ValidateRecord(record); err != nil {
		writeJSON(w, errorMsg{Error: err.Error()})
		return
	}

	// Remove the password just in case it marshals.
	record.Password = ""

	if err := i.db.Insert(record); err != nil {
		writeJSON(w, errorMsg{Error: "Could not create new user DB error"})
		log.Println(err)
		return
	}

	writeJSON(w, tokenMsg{Token: i.CreateToken(record)})
}

func (i *IDGen) TokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSON(w, errorMsg{Error: "Could not parse form data"})
		return
	}

	user := mux.Vars(r)["username"]
	pass := r.FormValue("password")

	record, err := i.FetchRecord(user, pass)

	if err != nil {
		writeJSON(w, errorMsg{Error: err.Error()})
		return
	}

	writeJSON(w, tokenMsg{Token: i.CreateToken(record)})
}
