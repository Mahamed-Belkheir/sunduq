package sunduq

import (
	"bufio"
	"reflect"
	"testing"
)

func assert(expected, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected: %v \n got: %v", expected, got)
	}
}

func serializeAndTest(msg Message, t *testing.T) Message {
	buf := msg.ToBytesBuffer()
	reader := bufio.NewReader(&buf)
	msg2, err := MessageFromBytes(reader)
	if err != nil {
		t.Errorf("failed to deserialize the message with error: %v", err)
	}
	assert(msg, msg2, t)
	return msg2
}

func TestPingSerialize(t *testing.T) {
	ping := NewPing(2)

	serializeAndTest(ping, t)
}

func TestResultSerialize(t *testing.T) {
	text := "task failed successfully"
	message := NewResult(1, true, String, []byte(text))
	message2 := serializeAndTest(message, t)

	if text != string(message2.Value) {
		t.Errorf("serialized message did not match the original, \n expected: %v \n got: %v", message, string(message2.Value))
	}
}

func TestConnectMessage(t *testing.T) {
	username := "bob123"
	password := "securepassword3000"
	con := NewConnect(username, password)
	con2 := serializeAndTest(con, t)

	assert(username, con2.Key, t)
	assert(password, string(con2.Value), t)
}

func TestDisconnectMessage(t *testing.T) {
	dsc := NewDisconnect()
	serializeAndTest(dsc, t)
}

func TestMessage(t *testing.T) {
	msg := NewMessage(Get, 1, "tuesday", "posts")
	serializeAndTest(msg, t)
}

func TestMessageWithValue(t *testing.T) {
	msg := NewMessageWithValue(Get, 1, "tuesday", "posts", String, []byte("some tuesday posts"))
	serializeAndTest(msg, t)
}
