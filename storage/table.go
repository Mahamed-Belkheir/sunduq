package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Mahamed-Belkheir/sunduq"
)

//AccessLevel is the base enum type for access lists
type AccessLevel uint8

const (
	//Read access allows reading from tables
	Read AccessLevel = iota
	//Write access allows reading and writing to tables
	Write
	//Admin access allows the above + editing access list and deleting tables
	Admin
)

func toAccessLevel(slice []byte) AccessLevel {
	return AccessLevel(slice[0])
}

//AccessList is a map of all users with access to the table and their access level
type AccessList map[string]AccessLevel

var accessMap map[sunduq.MessageType]AccessLevel = map[sunduq.MessageType]AccessLevel{
	sunduq.Get:          Read,
	sunduq.Set:          Write,
	sunduq.Del:          Write,
	sunduq.DeleteTable:  Admin,
	sunduq.SetTableUser: Admin,
	sunduq.DelTableUser: Admin,
}

func (t Table) verifyAccess(user string, command sunduq.MessageType) error {
	if command == sunduq.CreateTable {
		return nil
	}
	access, ok := t.ac[user]
	if !ok {
		return fmt.Errorf("user %v does not have access to table %v", user, t.name)
	}
	requiredLevel, ok := accessMap[command]
	if !ok {
		return fmt.Errorf("recieved invalid command id: %v", command)
	}
	if access < requiredLevel {
		return fmt.Errorf("user %v lacks the required access level: %v", user, requiredLevel)
	}
	return nil
}

//Table is the basic building block for the storage, holds all the data and includes other meta data
type Table struct {
	ac            AccessList
	data          map[string][]byte
	isSystemTable bool
	name          string
	mut           *sync.RWMutex
}

//NewTable creates a new Table instance
func NewTable(name, owner string, isSystemTable bool) Table {
	ac := make(AccessList)
	ac[owner] = Admin
	return Table{
		ac,
		make(map[string][]byte),
		isSystemTable,
		name,
		&sync.RWMutex{},
	}
}

func invalidCreds() error {
	return errors.New("User lacks the required access level")
}

func notFound() error {
	return errors.New("Key not found")
}

//Set adds an element to the table
func (t Table) Set(key string, value []byte) {
	t.mut.Lock()
	defer t.mut.Unlock()
	t.data[key] = value
}

//Get fetches the item from the table
func (t Table) Get(key string) ([]byte, error) {
	t.mut.RLock()
	defer t.mut.RUnlock()

	data, ok := t.data[key]
	if !ok {
		return nil, notFound()
	}

	return data, nil
}

//Del deletes the item from the table
func (t Table) Del(key string) {
	t.mut.Lock()
	defer t.mut.Unlock()
	delete(t.data, key)
}

//SetUser sets the user table access level
func (t Table) SetUser(user string, access AccessLevel) {
	t.mut.Lock()
	defer t.mut.Unlock()
	t.ac[user] = access
}

//DelUser deletes the user from the table access map
func (t Table) DelUser(user string) {
	t.mut.Lock()
	defer t.mut.Unlock()
	delete(t.ac, user)
}
