package event

import "github.com/Mahamed-Belkheir/sunduq"

//Manager handles all messaging events, and message middleware, connections and controllers register here
type Manager struct {
	index             handlerIndex
	recieveMiddleware []func(sunduq.Envelope) *sunduq.Envelope
	sendMiddleware    []func(sunduq.Envelope) *sunduq.Envelope
	recieveHandler    func(sunduq.Envelope)
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
	for _, handler := range e.recieveMiddleware {
		result := handler(envelope)
		if result == nil {
			return
		}
		envelope = *result
	}
	e.recieveHandler(envelope)
}

//Send passes the message from the send middleware and then into the proper connection
func (e Manager) Send(envelope sunduq.Envelope) {
	for _, handler := range e.sendMiddleware {
		result := handler(envelope)
		if result == nil {
			return
		}
		envelope = *result
	}
	e.index.send(envelope)

}

//RecieveHandler registers the function to manage recieved messages
func (e Manager) RecieveHandler(handler func(sunduq.Envelope)) {
	e.recieveHandler = handler
}

//UseOnRecieve add middleware functions to be run on any incoming message
func (e Manager) UseOnRecieve(handler func(sunduq.Envelope) *sunduq.Envelope) {
	e.recieveMiddleware = append(e.recieveMiddleware, handler)
}

//UseOnSend add middleware functions to be run on any outgoing message
func (e Manager) UseOnSend(handler func(sunduq.Envelope) *sunduq.Envelope) {
	e.sendMiddleware = append(e.sendMiddleware, handler)
}
