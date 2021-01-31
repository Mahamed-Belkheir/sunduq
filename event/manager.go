package event

import "github.com/Mahamed-Belkheir/sunduq"

//Manager handles all messaging events, and message middleware, connections and controllers register here
type Manager struct {
	index          *handlerIndex
	recieveHandler func(sunduq.Envelope)
}

//NewManager creates a new manager instance
func NewManager() Manager {
	return Manager{
		newHandlerIndex(),
		func(e sunduq.Envelope) {
			panic("no recieve handler is registered")
		},
	}
}

//RegisterConnection registers the connection to send the response to
func (e Manager) RegisterConnection(id int, channel chan sunduq.Message) {
	e.index.register(id, channel)
}

//UnregisterConnection removes the connection from the index, use before shutting down the connection
func (e Manager) UnregisterConnection(id int) {
	e.index.unregister(id)
}

//Recieve passes recieved messages through the middleware and into the controller
func (e Manager) Recieve(envelope sunduq.Envelope) {
	e.recieveHandler(envelope)
}

//Send passes the message from the send middleware and then into the proper connection
func (e Manager) Send(envelope sunduq.Envelope) {
	e.index.send(envelope)
}

//RecieveHandler registers the function to manage recieved messages
func (e *Manager) RecieveHandler(handler func(sunduq.Envelope)) {
	e.recieveHandler = handler
}
