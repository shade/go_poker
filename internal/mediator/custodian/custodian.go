// Custodian is for creating, validating, and removing
// tables.
package custodian

import (
	"go_poker/internal/cache"
	. "go_poker/internal/interfaces"
	msgpb "go_poker/internal/proto"
)

type Custodian struct {
	c cache.ICache
}

func NewCustodian(c ICache) {
	return Custodian{c}
}

func (c *Custodian) CreateTable() {
	opts := msgpb.TableOptions{}
}
