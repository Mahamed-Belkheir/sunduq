package storage

import (
	"fmt"
	"testing"

	"github.com/Mahamed-Belkheir/sunduq"
)

func TestStorageCreation(t *testing.T) {
	s := NewStorage()
	assert(s.tables, map[string]Table{}, t)
}

func TestStorageAddTable(t *testing.T) {
	s := NewStorage()
	user, table := "userOne", "tableOne"
	s.addTable(user, table, false)
	assert(Admin, s.tables[table].ac[user], t)
	assert(table, s.tables[table].name, t)
}

func TestStorageQuery(t *testing.T) {
	s := NewStorage()
	user, table, key, value := "userOne", "tableOne", "key1", []byte("hello world")
	s.addTable(user, table, false)

	_, err := s.
		Query(sunduq.Set).
		Table(table).
		User(user).
		Key(key).
		Value(value).
		Exec()

	ok(err, t)

	data, err := s.
		Query(sunduq.Get).
		Table(table).
		User(user).
		Key(key).
		Exec()

	ok(err, t)
	assert(value, data, t)

	_, err = s.
		Query(sunduq.Del).
		Table(table).
		User(user).
		Key(key).
		Exec()

	ok(err, t)

	data, err = s.
		Query(sunduq.Get).
		Table(table).
		User(user).
		Key(key).
		Exec()

	assert(notFound(), err, t)

}

func TestStorageTableQuery(t *testing.T) {
	s := NewStorage()
	user, table := "userOne", "tableOne"

	_, err := s.Query(sunduq.CreateTable).Table(table).User(user).Exec()
	ok(err, t)

	assert(NewTable(table, user, false), s.tables[table], t)

	newUser := "userTwo"
	_, err = s.Query(sunduq.SetTableUser).Table(table).User(user).Key(newUser).Value([]byte{uint8(Write)}).Exec()
	ok(err, t)

	_, err = s.Query(sunduq.Set).Table(table).User(newUser).Key("key").Value([]byte("hello world")).Exec()
	ok(err, t)

	_, err = s.Query(sunduq.SetTableUser).Table(table).User(newUser).Key(newUser).Value([]byte{uint8(Admin)}).Exec()
	assert(fmt.Errorf("user %v lacks the required access level: %v", newUser, Admin), err, t)

	_, err = s.Query(sunduq.DelTableUser).Table(table).User(user).Key(newUser).Exec()
	ok(err, t)

	_, err = s.Query(sunduq.Get).Table(table).User(newUser).Key("key").Exec()
	assert(fmt.Errorf("user %v does not have access to table %v", newUser, table), err, t)

	_, err = s.Query(sunduq.DeleteTable).Table(table).User(user).Exec()
	ok(err, t)

	_, err = s.Query(sunduq.Set).Table(table).User(user).Key("key").Value([]byte("hello world")).Exec()
	assert(fmt.Errorf("table %v not found", table), err, t)

}
