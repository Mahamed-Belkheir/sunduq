package storage

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(expected, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nexpected: %v \n got: %v", expected, got)
	}
}

func ok(err error, t *testing.T) {
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestCanCreateTable(t *testing.T) {
	tableName, user := "table1", "bob"
	table := NewTable(tableName, user, false)
	assert(table.name, tableName, t)
	assert(table.ac, AccessList{
		user: Admin,
	}, t)
	assert(table.data, map[string][]byte{}, t)
}

func TestTableCRUD(t *testing.T) {
	tableName, user := "table1", "bob"
	table := NewTable(tableName, user, false)

	value := []byte("world")
	table.Set("hello", value)
	result, err := table.Get("hello")
	ok(err, t)
	assert(result, value, t)
	table.Del("hello")
	_, err = table.Get("hello")
	assert(err, notFound(), t)
}

func TestAdminTableAccess(t *testing.T) {
	tableName, user := "table1", "bob"
	table := NewTable(tableName, user, false)

	for command := range accessMap {
		err := table.verifyAccess(user, command)
		ok(err, t)
	}
}

func TestWriteTableAccess(t *testing.T) {
	tableName, user := "table1", "bob"
	table := NewTable(tableName, user, false)

	newUser := "adam"

	table.SetUser(newUser, Write)
	for command, level := range accessMap {
		err := table.verifyAccess(newUser, command)
		if level == Admin {
			assert(fmt.Errorf("user %v lacks the required access level: %v", newUser, level), err, t)
			continue
		}
		ok(err, t)
	}
}

func TestReadTableAccess(t *testing.T) {
	tableName, user := "table1", "bob"
	table := NewTable(tableName, user, false)

	newUser := "adam"

	table.SetUser(newUser, Read)
	for command, level := range accessMap {
		err := table.verifyAccess(newUser, command)
		if level == Read {
			ok(err, t)
			continue
		}
		assert(fmt.Errorf("user %v lacks the required access level: %v", newUser, level), err, t)
	}
}

func TestNoTableAccess(t *testing.T) {
	tableName, user := "table1", "bob"
	table := NewTable(tableName, user, false)

	newUser := "john"

	for command := range accessMap {
		err := table.verifyAccess(newUser, command)
		assert(fmt.Errorf("user %v does not have access to table %v", newUser, tableName), err, t)
	}
}
