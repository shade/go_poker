package identity

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go_poker/internal/identity/db"
	"go_poker/internal/server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// 1 day in nanoseconds
const TOKEN_DURATION = 86486400000000000
const AUTH_HEADER = "Authorization"
const BEARER_SCHEMA = "Bearer "

type IDGen struct {
	db     db.IIDB
	secret string
	algo   jwt.SigningMethod
}

func NewIDGen(db db.IIDB, secret string) *IDGen {
	return &IDGen{
		db:     db,
		secret: secret,
		algo:   jwt.SigningMethodHS256,
	}
}

func (i *IDGen) IsValidPassword(hash string, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func (i *IDGen) FetchRecord(username string, password string) (*db.Record, error) {
	record, err := i.db.Get(db.DBKey(username))

	if err != nil {
		return nil, errors.New("User doesn't exist")
	}

	if !i.IsValidPassword(record.PasswordHash, password) {
		return nil, errors.New("Invalid password")
	}

	return record, nil
}

func (i *IDGen) ValidateRecord(r *db.Record) error {
	// Ensure all the fields are valid
	if len(r.Username) == 0 {
		return errors.New("Unset field: Username")
	}

	if len(r.Password) == 0 {
		return errors.New("Unset field: Password")
	}
	// Ensure no user with same username in the db
	record, err := i.db.Get(db.DBKey(r.Username))

	if record != nil {
		return errors.New("User already exists in DB")
	}

	if err != nil {
		return err
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

func (i *IDGen) CreateToken(r *db.Record) string {
	token, err := jwt.NewWithClaims(i.algo, jwt.MapClaims{
		"username": r.Username,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	}).SignedString([]byte(i.secret))

	if err != nil {
		log.Fatal(err)
	}

	return token
}

type tokenMsg struct {
	Token string `json:"token"`
}

func (i *IDGen) IDHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteJSON(w, utils.ErrorMsg{Error: "Could not parse form data"}, http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	user := r.FormValue("username")
	pass := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), 10)

	record := &db.Record{name, user, pass, base64.StdEncoding.EncodeToString(hash)}
	if err := i.ValidateRecord(record); err != nil {
		utils.WriteJSON(w, utils.ErrorMsg{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	// Remove the password just in case it marshals.
	record.Password = ""

	if err := i.db.Insert(record); err != nil {
		utils.WriteJSON(w, utils.ErrorMsg{Error: "Could not create new user DB error"}, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	utils.WriteJSON(w, tokenMsg{Token: i.CreateToken(record)}, http.StatusOK)
}

func (i *IDGen) TokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteJSON(w, utils.ErrorMsg{Error: "Could not parse form data"}, http.StatusBadRequest)
		return
	}

	user := mux.Vars(r)["username"]
	pass := r.FormValue("password")

	record, err := i.FetchRecord(user, pass)

	if err != nil {
		utils.WriteJSON(w, utils.ErrorMsg{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, tokenMsg{Token: i.CreateToken(record)}, http.StatusOK)
}

func (i *IDGen) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		splitToken := strings.Split(r.Header.Get(AUTH_HEADER), BEARER_SCHEMA)
		if len(splitToken) != 2 {
			utils.WriteJSON(w, utils.ErrorMsg{Error: "Invalid bearer token"}, http.StatusForbidden)
			return
		}

		token := strings.TrimSpace(splitToken[1])
		fmt.Println(token)
		valid, _ := i.ParseToken(token)
		if !valid {
			utils.WriteJSON(w, utils.ErrorMsg{Error: "Invalid bearer token"}, http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}
