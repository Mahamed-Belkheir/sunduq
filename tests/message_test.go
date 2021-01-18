package tests

import (
	"bufio"
	"reflect"
	"testing"

	"github.com/Mahamed-Belkheir/sunduq"
)

func TestPingSerialize(t *testing.T) {
	ping := sunduq.NewPing(2)
	buf := ping.ToBytesBuffer()
	reader := bufio.NewReader(&buf)
	rePing, err := sunduq.MessageFromBytes(reader)
	if err != nil {
		t.Errorf("ping deserialization failed with error %v", err)
	}

	if !reflect.DeepEqual(ping, rePing) {
		t.Errorf("serialized ping did not equal original ping message, \n ping: %v \n reping: %v", ping, rePing)
	}
}

func TestResultSerialize(t *testing.T) {
	message := "task failed successfully"
	res := sunduq.NewResult(1, true, sunduq.String, []byte(message))
	buf := res.ToBytesBuffer()
	reader := bufio.NewReader(&buf)
	reRes, err := sunduq.MessageFromBytes(reader)

	if err != nil {
		t.Errorf("result deserialization failed with error %v", err)
	}

	if !reflect.DeepEqual(res, reRes) {
		t.Errorf("serialized result did not equal original result message, \n res: %v \n reRes: %v", res, reRes)
	}

	if message != string(reRes.Value) {
		t.Errorf("serialized message did not match the original, \n expected: %v \n got: %v", message, string(reRes.Value))
	}
}

func TestConnectMessage(t *testing.T) {
	username := "bob123"
	password := "securepassword3000"
	con := sunduq.NewConnect(username, password)
	buf := con.ToBytesBuffer()
	reader := bufio.NewReader(&buf)
	reCon, err := sunduq.MessageFromBytes(reader)

	if err != nil {
		t.Errorf("failed to parse message with error: %v", err)
	}

	if !reflect.DeepEqual(con, reCon) {
		t.Errorf("original message did not match parsed message \n original: %v \n parsed: %v", con, reCon)
	}

}
