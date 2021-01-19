package sunduq

import (
	"bufio"
	"bytes"
	"fmt"
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

var sliceOfTypes = []MessageType{
	Ping,
	Result,
	Connect,
	Disconnect,
	Get,
	Set,
	Del,
	All,
	CreateTable,
	GetTables,
	DeleteTable,
	SetTableUser,
	DelTableUser,
	AllTableUser,
}
var sliceOfValueTypes = []ValueType{
	Boolean,
	String,
	Integer,
	Float,
	JSON,
}

func generateAllPossibleMessageCombinations() []Message {
	var i uint16
	i = 1
	messages := make([]Message, 0)
	for _, mType := range sliceOfTypes {
		for _, vType := range sliceOfValueTypes {
			isErr := false
			if i%2 == 0 {
				isErr = true
			}
			msg := Message{
				Error:     isErr,
				ID:        i,
				Key:       "msgKey " + fmt.Sprint(i),
				Table:     "some table",
				Type:      mType,
				ValueType: vType,
				Value:     []byte("message id is: " + fmt.Sprint(i)),
			}
			i++
			messages = append(messages, msg)
		}
	}
	return messages
}

func TestAllMessageCombinators(t *testing.T) {
	allMessages := generateAllPossibleMessageCombinations()
	for _, msg := range allMessages {
		serializeAndTest(msg, t)
	}
}

func TestParsingMessagesInARow(t *testing.T) {
	messages := generateAllPossibleMessageCombinations()
	var buf = bytes.NewBuffer([]byte{})
	for _, msg := range messages {
		msgBuf := msg.ToBytesBuffer()
		msgBuf.WriteTo(buf)
	}
	reader := bufio.NewReader(buf)
	serializedMessages := make([]Message, len(messages))
	lim := len(messages)
	var i int
	for i = 0; i < lim; i++ {
		msg, err := MessageFromBytes(reader)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Errorf("reading failed with non EOF error: %v", err)
		}
		serializedMessages[i] = msg
	}
	if len(messages) != len(serializedMessages) {
		t.Errorf("did not match the amount of messages serialized, got: %v, expected: %v", len(messages), len(serializedMessages))
	}
	for i, msg := range messages {
		assert(msg, serializedMessages[i], t)
	}
}
