package sunduq

//Handler connection handling interface, run method is ran by the server workers to listen to and send messages over a connection
type Handler interface {
	Run()
}

//Envelope is a struct that holds the message and the connection ID together
type Envelope struct {
	ID      int
	Message Message
	User    string
}

//NewEnvelope creates a new envelope struct
func NewEnvelope(id int, msg Message, user string) Envelope {
	return Envelope{
		id,
		msg,
		user,
	}
}
