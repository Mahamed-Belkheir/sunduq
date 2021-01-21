package tcp

import (
	"net"
	"reflect"
	"testing"

	"github.com/Mahamed-Belkheir/sunduq"
)

func assert(expected, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected: %v \n got: %v", expected, got)
	}
}

func TestConnection(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	server, client := NewConnection(serverConn), NewConnection(clientConn)
	server.Run()
	client.Run()
	msg := sunduq.NewPing(1)
	client.Send(msg)

	ch := server.Recieve()
	msg2 := <-ch

	assert(msg, msg2, t)
	client.Close()
	server.Close()
}
