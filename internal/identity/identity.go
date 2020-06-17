package identity

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

const PAYLOAD_MAX_SIZE = 100000

type Payload map[string][]byte


type IdentityGenerator struct {
	db *IdentityDB
	secret string
}

func mapSize(p Payload) int64 {
	i := int64(0)

	for k, v := range p {
		i += int64(len(k))
		i += int64(len(v))
	}

	return i
}

// NewIdentityGenerator creates an identity generator provided
// a specific token secret for creating HMAC keys.
func NewIdentityGenerator(secret string, db_url string) *IdentityGenerator {
	return &IdentityGenerator{
		secret,
	}
}

func (i *IdentityGenerator) NewIdentity(data Payload, unique string) ([]byte, error) {
	// Prevent memory exhaustion attacks
	if mapSize(data) > PAYLOAD_MAX_SIZE {
		return nil, errors.New("Provided data payload too large")
	}

	i.
	return nil, []byte(i.secret)
}

// NewToken Creates a new token provided a unique id,
// note this method is not idempotent
func (i *IdentityGenerator) NewToken(id string) []byte {

	return []byte(i.secret)
}

func (i *IdentityGenerator) Validate(t []byte) bool {
	// TOD
	return true
}

// TODO, refactor (move this to its own server package
func (i *IdentityGenerator) creationHandler(w http.ResponseWriter, r *http.Request) {

}

func (i *IdentityGenerator) tokenHandler(w http.ResponseWriter, r *http.Request) {

}

func (i *IdentityGenerator) CreateRouter(secret string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/gen", i.creationHandler)
	r.HandleFunc("/token", i.tokenHandler)

	return r
}
