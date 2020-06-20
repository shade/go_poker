package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	testRecord := &Record{
		Name:     "John",
		Username: "john",
		Password: "pass",
	}
	filedb := NewFileDB("../../test/test_db.csv")
	id := NewIDGen(filedb, "lol")
	token := id.CreateToken(testRecord)
	valid, _ := id.ParseToken(token)
	assert.True(t, valid)

	// Test the fail case
	token += "failpls"
	valid, _ = id.ParseToken(token)
	assert.False(t, valid)
}
