package sock

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

	msgpb "go_poker/internal/proto"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestPingPong(t *testing.T) {
	host := ":8080"
	// Run the server and don't block
	sock := NewSock()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got conn?, Upgrading.")
		sock.AddConnection(w, r)
	})
	go func() {
		err := http.ListenAndServe(host, nil)
		assert.Nil(t, err)
	}()

	// Wait 1 second for server to set up
	timer := time.NewTimer(1e9)
	<-timer.C
	// Attempt a ping pong message
	url := url.URL{Scheme: "ws", Host: host, Path: "/", RawQuery: "format=proto"}
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	assert.Nil(t, err)

	packet := msgpb.Packet{
		Event: msgpb.EventType_PING,
	}
	data, err := proto.Marshal(&packet)
	assert.Nil(t, err)
	err = c.WriteMessage(websocket.TextMessage, data)
	assert.Nil(t, err)

	// Attempt to read in the message
	dataChan := make(chan []byte)
	go func() {
		_, data, err = c.ReadMessage()
		assert.Nil(t, err)
		dataChan <- data
	}()

	// Max 2 seconds for the data to come back
	timer.Reset(1e9)
	select {
	case data = <-dataChan:
		proto.Unmarshal(data, &packet)
		assert.Equal(t, msgpb.EventType_PONG, packet.Event, "Expected PONG, got something else")
		return
	case <-timer.C:
		assert.Fail(t, "Test timeout, no PONG received")
		return
	}
}
