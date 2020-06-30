// Guides the new sockets to their respective rooms
package concierge

import (
	"gopoker/internal/identity"
	. "gopoker/internal/interfaces"
	"gopoker/internal/mediator/cache"
	msgpb "gopoker/internal/proto"
	"gopoker/internal/room"
	"gopoker/internal/room/user/sock"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
)

type Concierge struct {
	_IDGen  *identity.IDGen
	roomMap map[string]*room.Room
	sockMap map[string]ISock
	cache   cache.ICache
}

func NewConcierge(IDGen *identity.IDGen, cache cache.ICache) *Concierge {
	c := &Concierge{
		_IDGen:  IDGen,
		roomMap: map[string]*room.Room{},
		cache:   cache,
	}

	return c
}

func (c *Concierge) Start() {
	tblChan := make(chan string)
	c.cache.Poll(tblChan)

	for {
		value := <-tblChan
		// Potential timeout
		if value == "" {
			continue
		}
		tblOpts := msgpb.TableOptions{}
		proto.Unmarshal([]byte(value), &tblOpts)
		c.createRoom(&tblOpts)
	}
}

func (c *Concierge) createRoom(opts *msgpb.TableOptions) {
	// Create the room and update the room map
	room := room.NewRoom(opts)

	if _, ok := c.roomMap[opts.GetName()]; ok {
		// TODO: handle naming collisions by graceful update
		// of redis instance
	}

	c.roomMap[opts.GetName()] = room
}

func (c *Concierge) Resync() {
	// Resynchronize with the redis instance
	// If we can't connect to the redis instance, die.
}

func (c *Concierge) HandleConnection(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get(identity.AUTH_HEADER)

	if !strings.HasPrefix(auth, identity.BEARER_SCHEMA) {
		// TODO handle error.

		return
	}

	token := auth[len(identity.BEARER_SCHEMA):]

	// Find the user responsible for this socket
	userIds, ok := r.URL.Query()["userId"]

	if !ok || len(userIds) == 0 {
		// TODO, instant end.
	}

	splitToken := strings.Split(r.Header.Get(identity.AUTH_HEADER), identity.BEARER_SCHEMA)
	if len(splitToken) != 2 {
		// TODO: instant end.
	}

	token = strings.TrimSpace(splitToken[1])

	// Validate their token
	valid, claims := c._IDGen.ParseToken(token)

	if !valid {
		// TODO: instant end.
	}

	username, exists := claims["username"]
	if exists {
		// TODO: instant end.
		// Log, because this is a server issue.
	}

	s, exists := c.sockMap[username.(string)]

	if !exists {
		s = sock.NewSock()
		c.sockMap[username.(string)] = s
	}

	s.AddConnection(w, r)

	// BIG TODO!
	// If in room, grab the user and update their sock
	// Otherwise, make a new one and add them to the room
}
