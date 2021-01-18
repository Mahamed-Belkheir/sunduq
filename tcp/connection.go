package tcp

import (
	"bufio"
	"io"
	"net"

	"github.com/Mahamed-Belkheir/sunduq"
)

//Connection is an abstraction to the TCP Connection, handles serializing and sending, and recieving and parsing messages
type Connection struct {
	conn         net.Conn
	sendQueue    chan sunduq.Message
	recieveQueue chan sunduq.Message
}

//Recieve returns the recieved messages queue
func (h Connection) Recieve() chan sunduq.Message {
	return h.recieveQueue
}

//Send adds the message to the send queue to be sent over the connection by the Connection
func (h Connection) Send(msg sunduq.Message) {
	h.sendQueue <- msg
}

//Run runs the Connection to listen for new messages to recieve and send
func (h Connection) Run() {
	go func() {
		for {
			msg, ok := <-h.sendQueue
			if !ok {
				return
			}
			data := msg.ToBytesBuffer()
			io.Copy(h.conn, &data)
		}
	}()
	go func() {
		buf := bufio.NewReader(h.conn)
		auth, err := sunduq.MessageFromBytes(buf)
		_ = err
		// authorize
		_ = auth
		for {
			msg, err := sunduq.MessageFromBytes(buf)
			_ = err
			h.recieveQueue <- msg
		}
	}()
}
