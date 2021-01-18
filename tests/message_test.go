package tests

import (
	"bufio"
	"fmt"
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
		fmt.Println(ping, rePing)
		t.Errorf("serialized ping did not equal original ping message, \n ping: %v \n reping: %v", ping, rePing)
	}

}
