package sunduq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

//MessageType defines the message's type
type MessageType uint8

const (
	//Ping network test message
	Ping MessageType = iota
	//Result response to a message
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
	//CreateTable creates a new table
	CreateTable
	//DeleteTable deletes a table
	DeleteTable
	//SetTableUser sets table's user access list values
	SetTableUser
	//DelTableUser removes a user from a table access list
	DelTableUser
)

//ValueType the stored value type
type ValueType uint8

const (
	//Boolean ...
	Boolean ValueType = iota
	//String ...
	String
	//Integer ...
	Integer
	//Float ...
	Float
	//JSON ...
	JSON
	//Blob binary data
	Blob
)

/*
			  Message Format:

byte length:				purpose:

1			 		: 		MessageType
1			 		: 		Error (0/1)
2			 		: 		request id
1 		     		: 		table length
table length 		: 		table name
1					:		key length
key length			:		key
1                   :       ValueType
4					:		value length
value length		:		value
*/

//Message is the base message type that all messages share
type Message struct {
	Type      MessageType
	Error     bool
	ID        uint16
	Table     string
	Key       string
	ValueType ValueType
	Value     []byte
}

//NewPing creates a new Ping message
func NewPing(id uint16) Message {
	return Message{
		Type:  Ping,
		Error: false,
		ID:    id,
	}
}

//NewResult creates a new Result message
func NewResult(id uint16, err bool, valueType ValueType, value []byte) Message {
	return Message{
		Type:      Result,
		Error:     err,
		ID:        id,
		ValueType: valueType,
		Value:     value,
	}
}

//NewConnect creates a new Connect request message
func NewConnect(username, password string) Message {
	return Message{
		Type:      Connect,
		Error:     false,
		ValueType: String,
		Key:       username,
		Value:     []byte(password),
	}
}

//NewDisconnect creates a new disconnect request message
func NewDisconnect() Message {
	return Message{
		Type:  Disconnect,
		Error: false,
	}
}

//NewMessage creates a new value-less message, for queries
func NewMessage(mType MessageType, id uint16, key, table string) Message {
	return Message{
		Type:  mType,
		Error: false,
		ID:    id,
		Key:   key,
		Table: table,
	}
}

//NewMessageWithValue creates a message with a value, for set commands
func NewMessageWithValue(mType MessageType, id uint16, key, table string, valueType ValueType, value []byte) Message {
	return Message{
		Type:      mType,
		Error:     false,
		ID:        id,
		Key:       key,
		Table:     table,
		ValueType: valueType,
		Value:     value,
	}
}

//ToBytesBuffer serializes the message into bytes inside of a buffer to be transported over the network
func (m Message) ToBytesBuffer() bytes.Buffer {
	var buf bytes.Buffer

	buf.WriteByte(byte(m.Type))

	if m.Error {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}

	// if m.ID > 65535 {
	// 	return buf, errors.New("ID must be an int16 within 0 and 65535")
	// }
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, m.ID)
	buf.Write(id)

	tLen := len(m.Table)
	// if tLen > 255 {
	// 	return buf, errors.New("table name can not be longer than 255 characters")
	// }
	buf.WriteByte(byte(tLen))

	buf.Write([]byte(m.Table))

	kLen := len(m.Key)
	// if kLen > 255 {
	// 	return buf, errors.New("key name can not be longer than 255 characters")
	// }
	buf.WriteByte(byte(kLen))

	buf.Write([]byte(m.Key))

	buf.WriteByte(byte(m.ValueType))

	dataLen := make([]byte, 4)
	// if len(m.Value) > 4294967295 {
	// 	return buf, errors.New("value size must be within 4MB")
	// }
	if m.Value != nil {
		binary.LittleEndian.PutUint32(dataLen, uint32(len(m.Value)))
		buf.Write(dataLen)
		buf.Write(m.Value)
	} else {
		binary.LittleEndian.PutUint32(dataLen, 0)
		buf.Write(dataLen)
	}

	return buf
}

//MessageFromBytes reads the buffer for a valid message encoding
func MessageFromBytes(buf *bufio.Reader) (Message, error) {
	msg := Message{}

	typeByte, err := buf.ReadByte()
	if err != nil {
		return msg, fmt.Errorf("unable to read message type: %w", err)
	}
	msg.Type = MessageType(typeByte)

	errorByte, err := buf.ReadByte()
	if err != nil {
		return msg, fmt.Errorf("unable to read message status: %w", err)
	}
	if errorByte == 0 {
		msg.Error = false
	} else if errorByte == 1 {
		msg.Error = true
	} else {
		return msg, fmt.Errorf("invalid value: %v, for message status", errorByte)
	}

	idBytes := make([]byte, 2)
	_, err = buf.Read(idBytes)
	if err != nil {
		return msg, fmt.Errorf("unable to read message id: %w", err)
	}
	msg.ID = binary.LittleEndian.Uint16(idBytes)

	tableLen, err := buf.ReadByte()
	if err != nil {
		return msg, fmt.Errorf("unable to read message table length: %w", err)
	}
	if tableLen > 0 {
		tableName := make([]byte, tableLen)
		_, err = buf.Read(tableName)
		if err != nil {
			return msg, fmt.Errorf("unable to read table name: %w", err)
		}
		msg.Table = string(tableName)
	}

	keyLen, err := buf.ReadByte()
	if err != nil {
		return msg, fmt.Errorf("unable to read message key length: %w", err)
	}
	if keyLen > 0 {
		key := make([]byte, keyLen)
		_, err = buf.Read(key)
		if err != nil {
			return msg, fmt.Errorf("unable to read message key: %w", err)
		}
		msg.Key = string(key)
	}

	valueType, err := buf.ReadByte()
	if err != nil {
		return msg, fmt.Errorf("unable to read message value type: %w", err)
	}
	msg.ValueType = ValueType(valueType)

	valueLenBytes := make([]byte, 4)
	_, err = buf.Read(valueLenBytes)
	if err != nil {
		return msg, fmt.Errorf("unable to message value length: %w", err)
	}
	valueLen := binary.LittleEndian.Uint32(valueLenBytes)
	if valueLen > 0 {
		value := make([]byte, valueLen)
		buf.Read(value)
		msg.Value = value
	}

	return msg, nil
}
