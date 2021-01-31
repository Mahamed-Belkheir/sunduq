package server

import (
	"errors"
	"net"
	"time"

	"github.com/Mahamed-Belkheir/sunduq"
	"github.com/Mahamed-Belkheir/sunduq/event"
	"github.com/Mahamed-Belkheir/sunduq/tcp"
)

//TCPServer is a TCP implementation for the server using the connection pool
type TCPServer struct {
}

//TCPHandler the TCP implementation for the Handler interface
type TCPHandler struct {
	conn   net.Conn
	id     int
	events *event.Manager
	auth   func(sunduq.Message) *sunduq.Message
	close  chan bool
}

//Run the function to start handling tcp connections
func (t TCPHandler) Run() {
	connection := tcp.NewConnection(t.conn)
	connection.Run()

	defer func() {
		t.events.UnregisterConnection(t.id)
		connection.Close()
		// log it
	}()

	t.events.RegisterConnection(t.id, connection.SendQueue)
	//todo get user via auth
	authRequest := <-connection.RecieveQueue
	if authRequest.Type != sunduq.Connect {
		//log it
		return
	}
	errResponse := t.auth(authRequest)
	if errResponse != nil {
		// log it
		connection.Send(*errResponse)
		time.Sleep(1) // give time for the response to be sent
		return
	}

	user := authRequest.Key
	errorQueue := connection.Errors()
	for {
		select {
		case err := <-errorQueue:
			//log it
			if errors.Unwrap(err).Error() == "EOF" {
				return
			}
		case msg := <-connection.RecieveQueue:
			t.events.Send(sunduq.NewEnvelope(t.id, msg, user))
		case <-t.close:
			return
		}
	}
}
