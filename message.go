package sunduq

//MessageType defines the message's type
type MessageType uint8

const (
	//Ping network test message
	Ping MessageType = iota
	//Response response to a nonquery reque
	Response
	//Result response to a query
	Result
	//Connect Establish a Connection
	Connect
	//Disconnect Close a connection
	Disconnect
	//Get a key value from a table
	Get
	//Set a key value on a table
	Set
	//Del delete a value from a table
	Del
	//All get all values from a table
	All
	//CreateTable creates a new table
	CreateTable
	//GetTables fetches all table names
	GetTables
	//DeleteTable deletes a table
	DeleteTable
	//SetTableUser sets table's user access list values
	SetTableUser
	//DelTableUser removes a user from a table access list
	DelTableUser
	//AllTableUser fetches all users in the access list
	AllTableUser
)

/*
Message Formats

Ping

1 : 0 // message start
1 : Ping // MessageType
1 : 0 // message end

Response

1 : 0 // message start
1 : Response // MessageType
1 : 0/1 // 0 = OK / 1 = ERROR, if OK ends here
4 : message length
message length: message
1 : 0 // message end


Result

1 : 0 // message start
1 : Result // MessageType
4 : // request Id
1 : // ResultType
4 : // result length
result length: // result
1 : 0 // message end


Connect

1 : 0 // message start
1 : Connect // MessageType
1 : // username length
username length: // username
1 : // password length
password length // password
1 : 0 // message end


Disconnect

1 : 0 // message start
1 : Disconnect // MessageType
1 : 0 // message end

Queries

1 : 0 // message start
1 : // MessageType
4 : // request Id
1 : // result type
4 : // data length
data length: // data
1 : 0 // message end



*/

//Message is the base message type that all messages share
type Message interface {
	Type() MessageType
}

//Handler is a generic connection handler interface, is supposed to manage handling sending and recieving messages from connections
type Handler interface {
	Recieve(interface{}) Message
	Send(Message)
	Run()
}
