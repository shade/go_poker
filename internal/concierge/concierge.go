// Guides the new sockets to their respective rooms
package concierge

import (
	"go_poker/internal/room"
	"go_poker/internal/queue"
	"net/http"
)

type Concierge struct {
	roomMap map[string]*room.Room
	queue queue.IQueue
}

func NewConcierge(queue queue.IQueue) *Concierge {
	c := &Concierge{
		roomMap: map[string]{},
		queue: queue,
	}

	c.Poll()
	return c
}

func (c *Concierge) Poll() {
	for {
		value, err := queue.Poll()

		// Potential timeout
		if value == nil {
			continue
		}

		// Potential server fail
		if err != nil {
			// TODO: handle queue failure case
			panic("Polling failed, queue server down")
		}

		c.CreateTable(value)
	}
}

func (c *Concierge) CreateRoom() {
	// Create the room and update the room map
}

func (c *Concierge) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Find the user responsible for this socket

	// Validate their token
	// If in room, grab the user and update their sock
	// Otherwise, make a new one and add them to the room
}
