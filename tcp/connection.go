package tcp

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/Mahamed-Belkheir/sunduq"
)

//Connection is an abstraction to the TCP Connection, handles serializing and sending, and recieving and parsing messages
type Connection struct {
	conn         net.Conn
	SendQueue    chan sunduq.Message
	RecieveQueue chan sunduq.Message
	errorQueue   chan error
	close        chan bool
	running      bool
	wg           *sync.WaitGroup
}

//NewConnection wraps the tcp connection with the Connection abstraction struct
func NewConnection(conn net.Conn) Connection {
	c := Connection{
		conn,
		make(chan sunduq.Message, 5),
		make(chan sunduq.Message, 5),
		make(chan error, 5),
		make(chan bool),
		false,
		&sync.WaitGroup{},
	}
	c.wg.Add(2)
	return c
}

//Recieve returns the recieved messages queue
func (h Connection) Recieve() chan sunduq.Message {
	return h.RecieveQueue
}

//Send adds the message to the send queue to be sent over the connection by the Connection
func (h Connection) Send(msg sunduq.Message) {
	h.SendQueue <- msg
}

//Errors returns the errors channel, errors recieved while the connection is processing messages
func (h Connection) Errors() chan error {
	return h.errorQueue
}

//Close closes the connection and ends all associated tasks
func (h Connection) Close() {
	h.conn.Close()
	close(h.close)
	close(h.SendQueue)
	h.wg.Wait()
	close(h.errorQueue)
}

//Run runs the Connection to listen for new messages to recieve and send
func (h Connection) Run() error {
	if h.running {
		return errors.New("connection is already running")
	}
	h.running = true
	go func() {
		defer h.wg.Done()
		for {
			msg, ok := <-h.SendQueue
			if !ok {
				return
			}
			data := msg.ToBytesBuffer()
			_, err := io.Copy(h.conn, &data)
			if err != nil {
				h.errorQueue <- err
				if errors.Unwrap(err).Error() == "EOF" {
					return
				}
			}
		}
	}()
	go func() {
		defer h.wg.Done()
		buf := bufio.NewReader(h.conn)
		for {
			select {
			case <-h.close:
				close(h.RecieveQueue)
				return
			default:
				msg, err := sunduq.MessageFromBytes(buf)
				if err != nil {
					h.errorQueue <- err
					if errors.Unwrap(err).Error() == "EOF" {
						close(h.RecieveQueue)
						return
					}
				} else {
					h.RecieveQueue <- msg
				}
			}
		}
	}()
	return nil
}
