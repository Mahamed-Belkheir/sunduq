package server

import (
	"errors"
	"net"

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
	close  chan bool
}

//Run the function to start handling tcp connections
func (t TCPHandler) Run() {
	connection := tcp.NewConnection(t.conn)
	defer func() {
		t.events.UnregisterConnection(t.id)
		connection.Close()
		// log it
	}()
	t.events.RegisterConnection(t.id, connection.SendQueue)
	errorQueue := connection.Errors()
	for {
		select {
		case err := <-errorQueue:
			//log it
			if errors.Unwrap(err).Error() == "EOF" {
				return
			}
		case msg := <-connection.RecieveQueue:
			t.events.Send(sunduq.NewEnvelope(t.id, msg))
		case <-t.close:
			return
		}
	}
}
