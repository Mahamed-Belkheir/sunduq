package tcp

import (
	"errors"
	"net"
	"reflect"
	"sync"
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

func TestRunAndClose(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	server, client := NewConnection(serverConn), NewConnection(clientConn)

	server.Run()
	client.Run()

	server.Close()
	client.Close()
}

func TestContinousMessaging(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	server, client := NewConnection(serverConn), NewConnection(clientConn)

	server.Run()
	client.Run()

	messages := []sunduq.Message{
		sunduq.NewPing(1),
		sunduq.NewPing(2),
		sunduq.NewPing(3),
		sunduq.NewPing(4),
		sunduq.NewPing(5),
		sunduq.NewPing(6),
		sunduq.NewPing(7),
		sunduq.NewPing(8),
		sunduq.NewPing(9),
		sunduq.NewPing(10),
	}
	wg := sync.WaitGroup{}
	wg.Add(10)

	go func() {
		for _, msg := range messages {
			client.Send(msg)
		}
	}()

	recievedMsgs := []sunduq.Message{}

	errRecieved := make(chan error)

	go func() {
		msgChan := server.Recieve()
		errChan := server.Errors()
		for {
			select {
			case msg := <-msgChan:
				recievedMsgs = append(recievedMsgs, msg)
				wg.Done()
			case err := <-errChan:
				errRecieved <- err
				return
			}
		}
	}()
	wg.Wait()
	client.Close()
	err := <-errRecieved
	server.Close()
	assert(errors.New("EOF"), errors.Unwrap(err), t)
	for i := range recievedMsgs {
		assert(messages[i], recievedMsgs[i], t)
	}
}
