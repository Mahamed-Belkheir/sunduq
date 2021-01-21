package tcp

import (
	"bufio"
	"errors"
	"io"
	"net"

	"github.com/Mahamed-Belkheir/sunduq"
)

//Connection is an abstraction to the TCP Connection, handles serializing and sending, and recieving and parsing messages
type Connection struct {
	conn         net.Conn
	sendQueue    chan sunduq.Message
	recieveQueue chan sunduq.Message
	errorQueue   chan error
	close        chan bool
	running      bool
}

//NewConnection wraps the tcp connection with the Connection abstraction struct
func NewConnection(conn net.Conn) Connection {
	return Connection{
		conn,
		make(chan sunduq.Message),
		make(chan sunduq.Message),
		make(chan error),
		make(chan bool),
		false,
	}
}

//Recieve returns the recieved messages queue
func (h Connection) Recieve() chan sunduq.Message {
	return h.recieveQueue
}

//Send adds the message to the send queue to be sent over the connection by the Connection
func (h Connection) Send(msg sunduq.Message) {
	h.sendQueue <- msg
}

//Errors returns the errors channel, errors recieved while the connection is processing messages
func (h Connection) Errors() chan error {
	return h.errorQueue
}

//Close closes the connection and ends all associated tasks
func (h Connection) Close() {
	h.conn.Close()
	close(h.close)
	close(h.sendQueue)
	close(h.errorQueue)
}

//Run runs the Connection to listen for new messages to recieve and send
func (h Connection) Run() error {
	if h.running {
		return errors.New("connection is already running")
	}
	h.running = true
	go func() {
		for {
			msg, ok := <-h.sendQueue
			if !ok {
				return
			}
			data := msg.ToBytesBuffer()
			_, err := io.Copy(h.conn, &data)
			if err != nil {
				h.errorQueue <- err
			}
		}
	}()
	go func() {
		buf := bufio.NewReader(h.conn)
		for {
			select {
			case <-h.close:
				close(h.recieveQueue)
				return
			default:
				msg, err := sunduq.MessageFromBytes(buf)
				if err != nil {
					h.errorQueue <- err
				}
				h.recieveQueue <- msg
			}
		}
	}()
	return nil
}
