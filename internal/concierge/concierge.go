// Guides the new sockets to their respective rooms
package concierge

import (
	"go_poker/internal/identity"
	"go_poker/internal/room"
	"go_poker/internal/cache"
	"net/http"
)

type Concierge struct {
	_IDGen *identity.IDGen 
	roomMap map[string]*room.Room
	cache cache.ICache
}

func NewConcierge(IDGen *identity.IDGen, cache cache.ICache) *Concierge {
	c := &Concierge{
		_IDGen: IDGen,
		roomMap: map[string]{},
		cache: cache,
	}

	go c.Poll()
	return c
}

func (c *Concierge) Poll() {
	for {
		value, err := cache.Poll()

		// Potential timeout
		if value == nil {
			continue
		}

		// Potential server fail
		if err != nil {
			// TODO: handle cache failure case
			panic("Polling failed, cache server down")
		}



		c.createRoom(value)
	}
}

func (c *Concierge) createRoom(opts *msgpb.TableOptions) {
	// Create the room and update the room map
	room := room.NewRoom(opts)

	if _, ok := c.roomMap[room.GetID()]; ok{
		// TODO: handle naming collisions by graceful update
		// of redis instance
	}

	c.roomMap[room.GetID()] = room
}

func (c *Concierge) Resync() {
	// Resynchronize with the redis instance
	// If we can't connect to the redis instance, die.
}

func (c *Concierge) HandleConnection(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get(AUTH_HEADER)

	if !strings.HasPrefix(auth, BEARER_SCHEMA) {
		// TODO handle error.

		return
	}

	token := auth[len(BEARER_SCHEMA):]
	
	// Find the user responsible for this socket
	userIds, ok := r.URL.Query()["userId"]

	if !ok || len(userIds) == 0 {
		// TODO, instant end.
	}
	
	splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(splitToken) != 2 {
		// TODO: instant end.
	}

	token = strings.TrimSpace(splitToken[1])

	// Validate their token
	// If in room, grab the user and update their sock
	// Otherwise, make a new one and add them to the room
}
