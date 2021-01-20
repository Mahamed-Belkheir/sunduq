package storage

import (
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
		Query().
		Table(table).
		User(user).
		Command(sunduq.Set).
		Key(key).
		Value(value).
		Exec()

	ok(err, t)

	data, err := s.
		Query().
		Table(table).
		User(user).
		Command(sunduq.Get).
		Key(key).
		Exec()

	ok(err, t)
	assert(value, data, t)

	_, err = s.
		Query().
		Table(table).
		User(user).
		Command(sunduq.Del).
		Key(key).
		Exec()

	ok(err, t)

	data, err = s.
		Query().
		Table(table).
		User(user).
		Command(sunduq.Get).
		Key(key).
		Exec()

	assert(notFound(), err, t)

}
