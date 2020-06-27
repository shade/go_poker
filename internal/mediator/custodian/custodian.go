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

func (c *Custodian) validateTable(opts *msgpb.TableOptions{}) error {
	// Check for duplicate keys
	if c.db.Get(opts.Name) != nil {
		return errors.New("Table name already taken.")
	}

	if opts.MinBuy <= 0 {
		return errors.New("Min buy must be greater than 0.")
	}

	if opts.MaxSeats < 2 {
		return errors.New("Your table must have at least 2 seats.")
	}

	if opts.BigBlind > opts.MinBuy {
		return errors.New("The big blind must be at least the minimum buyin.")
	}

	if math.MaxUint64 > (uint64(opts.MaxSeats) * opts.MinBuy) {
		return errors.New("The product of the max seats and min buy is too large.")
	}

	return nil
}

func (c *Custodian) CreateTable(w http.ResponseWriter, r *http.Request) {
	var opts msgpb.TableOptions

    err := json.NewDecoder(r.Body).Decode(&opts)
    if err != nil {
        utils.WriteJSON(w, utils.ErrorMsg{Error: err.Error()}, http.StatusBadRequest)
        return
	}
	
	err = c.validateTable(opts)
    if err != nil {
        utils.WriteJSON(w, utils.ErrorMsg{Error: err.Error()}, http.StatusBadRequest)
        return
	}

	err = c.createTable(opts)
    if err != nil {
        utils.WriteJSON(w, utils.ErrorMsg{Error: err.Error()}, http.StatusBadRequest)
        return
	}

	utils.WriteJSON(w, nil, http.StatusOK)
}

func (c *Custodian) FetchTables(w http.ResponseWriter, r *http.Request) {
	tables := c.getTables()

	utils.WriteJSON(w, tables, http.StatusOK)
}
