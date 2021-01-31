package event

import (
	"reflect"
	"testing"

	"github.com/Mahamed-Belkheir/sunduq"
)

func assert(expected, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nexpected: %v \n got: %v", expected, got)
	}
}

func TestHandlerIndex(t *testing.T) {
	handler := newHandlerIndex()
	ch1, ch2 := make(chan sunduq.Message), make(chan sunduq.Message)
	handler.register(1, ch1)
	handler.register(2, ch2)
	msg1, msg2 := sunduq.NewPing(1), sunduq.NewPing(2)
	go func() {
		handler.send(sunduq.NewEnvelope(1, msg1, "user"))
		handler.send(sunduq.NewEnvelope(2, msg2, "user"))
	}()

	rmsg1 := <-ch1
	rmsg2 := <-ch2

	assert(msg1, rmsg1, t)
	assert(msg2, rmsg2, t)
	handler.unregister(1)
	handler.unregister(2)
}
