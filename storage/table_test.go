package storage

import (
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
