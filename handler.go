package sunduq

//Handler connection handling interface, run method is ran by the server workers to listen to and send messages over a connection
type Handler interface {
	Run()
}
