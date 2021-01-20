package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Mahamed-Belkheir/sunduq"
)

func tableExists() error {
	return errors.New("Table already exists")
}

func tableNotFound() error {
	return errors.New("Table not found")
}

//Storage struct holds all the tables and holds all data manipulation methods
type Storage struct {
	tables map[string]Table
	mut    *sync.RWMutex
}

//NewStorage creates a new storage instance
func NewStorage() Storage {
	return Storage{
		make(map[string]Table, 0),
		&sync.RWMutex{},
	}
}

func (s Storage) addTable(user, tableName string, isSystemTable bool) error {
	defer s.mut.Unlock()
	s.mut.Lock()
	_, ok := s.tables[tableName]
	if ok {
		return tableExists()
	}
	s.tables[tableName] = NewTable(tableName, user, isSystemTable)
	return nil
}

func (s Storage) delTable(tableName string) error {
	defer s.mut.Unlock()
	s.mut.Lock()
	table, ok := s.tables[tableName]
	if !ok {
		return nil
	}
	if table.isSystemTable {

		return errors.New("Can not delete system table")
	}
	delete(s.tables, tableName)
	return nil
}

//Query returns a new query instance, can be chained to make a query to the storage
func (s Storage) Query(query sunduq.MessageType) QueryBuilder {
	q := QueryBuilder{}
	q.store = s
	q.command = query
	return q
}

//QueryBuilder exposes a fluent API to make queries to the storage tables
type QueryBuilder struct {
	store   Storage
	user    string
	table   string
	command sunduq.MessageType
	key     string
	value   []byte
}

//Table sets the query target table
func (q QueryBuilder) Table(table string) QueryBuilder {
	q.table = table
	return q
}

//User sets the user making the query, for authorization and not part of the query
func (q QueryBuilder) User(user string) QueryBuilder {
	q.user = user
	return q
}

//Key sets the query's target key
func (q QueryBuilder) Key(key string) QueryBuilder {
	q.key = key

	return q
}

//Value sets the query's Value, for Set* methods
func (q QueryBuilder) Value(value []byte) QueryBuilder {
	q.value = value
	return q
}

//Exec executes the query, returns a byte slice if it was a Get* request, and an error
func (q QueryBuilder) Exec() ([]byte, error) {
	if q.table == "" {
		return nil, errors.New("did not select a table for the query builder")
	}
	if q.user == "" {
		return nil, errors.New("did not select a user for the query builder")
	}
	if q.command == 0 {
		return nil, errors.New("did not select a command for the query builder")
	}

	q.store.mut.RLock()
	table, ok := q.store.tables[q.table]
	q.store.mut.RUnlock()
	if !ok {
		return nil, fmt.Errorf("table %v not found", q.table)
	}

	err := table.verifyAccess(q.user, q.command)
	if err != nil {
		return nil, err
	}

	switch q.command {
	case sunduq.Get:
		return table.Get(q.key)
	case sunduq.Set:
		if q.key == "" {
			return nil, errors.New("did not set a key for a Get query builder")
		}
		if q.value == nil {
			return nil, errors.New("did not set a value for a Set query builder")
		}
		table.Set(q.key, q.value)
		return nil, nil
	case sunduq.Del:
		if q.key == "" {
			return nil, errors.New("did not set a key for a Delete query builder")
		}
		table.Del(q.key)
		return nil, nil
	case sunduq.CreateTable:
		return nil, q.store.addTable(q.user, q.table, false)
	case sunduq.DeleteTable:
		return nil, q.store.delTable(q.table)
	case sunduq.SetTableUser:
		if q.key == "" {
			return nil, errors.New("did not set a key for a Get query builder")
		}
		if q.value == nil {
			return nil, errors.New("did not set a value for a Set query builder")
		}
		table.SetUser(q.key, toAccessLevel(q.value))
		return nil, nil
	case sunduq.DelTableUser:
		if q.key == "" {
			return nil, errors.New("did not set a key for a Get query builder")
		}
		table.DelUser(q.key)
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid query command ID recieved: %v", q.command)
	}
}
