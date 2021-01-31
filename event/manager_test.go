package event

import (
	"testing"

	"github.com/Mahamed-Belkheir/sunduq"
)

func TestManager(t *testing.T) {
	manager := NewManager()
	ch1 := make(chan sunduq.Message)

	manager.RegisterConnection(1, ch1)
	msg := sunduq.NewPing(1)

	go manager.Send(sunduq.NewEnvelope(1, msg, "user"))

	rmsg := <-ch1

	assert(msg, rmsg, t)

	recieved := false

	manager.RecieveHandler(func(env sunduq.Envelope) {
		recieved = true
	})

	manager.Recieve(sunduq.NewEnvelope(1, msg, "user"))

	assert(true, recieved, t)
}
