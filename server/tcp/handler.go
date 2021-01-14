package tcp

import (
	"net"

	"github.com/Mahamed-Belkheir/sunduq"
)

//Handler is the TCP implementation for the Handler interface, manages incoming and outgoing messages from a single TCP connection
type Handler struct {
	conn         net.Conn
	sendQueue    chan sunduq.Message
	recieveQueue chan sunduq.Message
}

//Recieve returns the recieved messages queue
func (h Handler) Recieve() chan sunduq.Message {
	return h.recieveQueue
}

//Send adds the message to the send queue to be sent over the connection by the Handler
func (h Handler) Send(msg sunduq.Message) {
	h.sendQueue <- msg
}

//Run runs the handler to listen for new messages to recieve and send
func (h Handler) Run() {
	go func() {
		for {
			// msg <- h.sendQueue
			// h.conn.Write()
		}
	}()
	go func() {
		for {
			// h.conn.Read()
			// h.recieveQueue <- msg
		}
	}()
}
