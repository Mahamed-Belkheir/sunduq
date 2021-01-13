package storage

import (
	"errors"
	"sync"

	"github.com/Mahamed-Belkheir/sunduq/context"
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

//IsValid verifies the access level is a valid value within the enum values
func (a AccessLevel) IsValid() error {
	switch a {
	case Read, Write, Admin:
		return nil
	default:
		return errors.New("Invalid AccessLevel Value")
	}
}

//AccessList is a map of all users with access to the table and their access level
type AccessList map[string]AccessLevel

func (a AccessList) canRead(user string) bool {
	_, ok := a[user]
	return ok
}

func (a AccessList) canWrite(user string) bool {
	access, ok := a[user]
	if !ok {
		return false
	}
	if access == Write {
		return true
	}
	return false
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

//Add adds an element to the table
func (t Table) Add(context context.Context, key string, value []byte) error {
	t.mut.Lock()
	defer t.mut.Unlock()
	if !t.ac.canWrite(context.User) {
		return invalidCreds()
	}

	t.data[key] = value

	return nil
}

//Get fetches the item from the table
func (t Table) Get(context context.Context, key string) ([]byte, error) {
	t.mut.RLock()
	defer t.mut.RUnlock()
	if !t.ac.canRead(context.User) {
		return nil, invalidCreds()
	}

	data, ok := t.data[key]
	if !ok {
		return nil, notFound()
	}

	return data, nil
}

//Del deletes the item from the table
func (t Table) Del(context context.Context, key string) error {
	t.mut.Lock()
	defer t.mut.Unlock()
	if !t.ac.canWrite(context.User) {
		return invalidCreds()
	}

	delete(t.data, key)

	return nil
}

//All fetches all data in the table
func (t Table) All(context context.Context, key string) (map[string][]byte, error) {
	t.mut.RLock()
	defer t.mut.RUnlock()
	if !t.ac.canRead(context.User) {
		return nil, invalidCreds()
	}

	return t.data, nil
}
