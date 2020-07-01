package room

import (
	msgpb "gopoker/internal/proto"
	"gopoker/internal/room/table/player"
	"gopoker/internal/room/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomCreation(t *testing.T) {
	NewRoom(&msgpb.TableOptions{
		Name:              "room1",
		Owner:             "joe",
		MinBuy:            0,
		MaxSeats:          10,
		SeatShuffle:       false,
		SeatShuffleRounds: 0,
		Rate:              0,
		BustedRebuyTime:   1000,
		RoundDelay:        1000,
		ActionTimeout:     10,
	})
}
func TestRoomSeatTaken(t *testing.T) {
	room := NewRoom(&msgpb.TableOptions{
		Name:              "room1",
		Owner:             "joe",
		MinBuy:            0,
		MaxSeats:          10,
		SeatShuffle:       false,
		SeatShuffleRounds: 0,
		Rate:              0,
		BustedRebuyTime:   1000,
		RoundDelay:        1000,
		ActionTimeout:     10,
	})

	p1 := player.NewPlayer(user.NewUser("u1"), 100, 1)
	p2 := player.NewPlayer(user.NewUser("u2"), 100, 1)

	assert.Nil(t, room.Table().SeatPlayer(p1), "Seat not taken?")
	assert.NotNil(t, room.Table().SeatPlayer(p2), "Seat taken!")
}

func TestRoomMinBuy(t *testing.T) {
	room := NewRoom(&msgpb.TableOptions{
		Name:              "room1",
		Owner:             "joe",
		MinBuy:            10,
		MaxSeats:          10,
		SeatShuffle:       false,
		SeatShuffleRounds: 0,
		Rate:              0,
		BustedRebuyTime:   1000,
		RoundDelay:        1000,
		ActionTimeout:     10,
	})

	p1 := player.NewPlayer(user.NewUser("u1"), 3, 1)
	p2 := player.NewPlayer(user.NewUser("u2"), 10, 1)

	assert.NotNil(t, room.Table().SeatPlayer(p1), "Bad buy")
	assert.Nil(t, room.Table().SeatPlayer(p2), "Good buy")
}

func TestRoomKick(t *testing.T) {
	room := NewRoom(&msgpb.TableOptions{
		Name:              "room1",
		Owner:             "joe",
		MinBuy:            10,
		MaxSeats:          10,
		SeatShuffle:       false,
		SeatShuffleRounds: 0,
		Rate:              0,
		BustedRebuyTime:   1000,
		RoundDelay:        1000,
		ActionTimeout:     10,
	})

	p1 := player.NewPlayer(user.NewUser("u1"), 3, 1)
	p2 := player.NewPlayer(user.NewUser("u2"), 10, 1)

	assert.NotNil(t, room.Table().SeatPlayer(p1), "Bad buy")
	assert.Nil(t, room.Table().SeatPlayer(p2), "Good buy")

}
