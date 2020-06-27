// Custodian is for creating, validating, and removing
// tables.
package custodian

import (
	"errors"
	. "go_poker/internal/interfaces"
	"go_poker/internal/mediator/cache"
	msgpb "go_poker/internal/proto"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type Custodian struct {
	db cache.ICache
}

func NewCustodian(c ICache) {
	return Custodian{c}
}

func (c *Custodian) createTable(opts *msgpb.TableOptions) error {
	optBytes, err := proto.Marshal(opts)

	if err != nil {
		return errors.New("Could not serialize options, invalid.")
	}

	c.db.Push(opts.Name, optBytes)
}

func (c *Custodian) getTables() proto.Message {
	c := []*proto.Message{}

	for key := range c.db.Keys() {
		msg := msgpb.TableOptions{}
		data := c.db.Get(key)
		proto.Unmarshal(data, msg)

		c = append(c, msg)
	}

	return c
}

func (c *Custodian) CreateTables(w http.ResponseWriter, r *http.Request) {
}

func (c *Custodian) FetchTables(w http.ResponseWriter, r *http.Request) {
	tables := c.getTables()

}
