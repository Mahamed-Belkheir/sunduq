package sunduq

//Request struct holds the request details
type Request struct {
	ID   uint32
	User string
}

//NewRequest creates a new request struct
func NewRequest(id uint32, user string) Request {
	return Request{
		id, user,
	}
}
