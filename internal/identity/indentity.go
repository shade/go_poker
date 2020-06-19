package identity

import (
	"errors"
	"net/http"
)

type IDGen struct {
	db     IIDB
	secret string
}

func NewIDGen(db IIDB, secret string) *IDGen {
	return &IDGen{
		db:     db,
		secret: secret,
	}
}

func (i *IDGen) ValidateRecord(r *Record) (bool, error) {
	// Ensure no user with same username in the db

	if record, _ := i.db.Get(DBKey(r.Username)); record != nil {
		return false, errors.New("User already exists in DB")
	}

	return true, nil
}

func (i *IDGen) IsValidToken(token string) bool {
	return true
}

func (i *IDGen) CreateToken(username string, password string) string {
	return ""
}

func (i *IDGen) TokenHandler(w http.ResponseWriter, r *http.Request) {

}
func (i *IDGen) IDHandler(w http.ResponseWriter, r *http.Request) {

}
