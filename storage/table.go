package storage

import (
	"errors"
)

//AccessLevel is the base enum type for access lists
type AccessLevel uint8

const (
	//Read access allows reading from tables
	Read AccessLevel = iota
	//Write access allows reading and writing to tables
	Write
)

//IsValid verifies the access level is a valid value within the enum values
func (a AccessLevel) IsValid() error {
	switch a {
	case Read, Write:
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
	ac   AccessList
	data map[string][]byte
}
