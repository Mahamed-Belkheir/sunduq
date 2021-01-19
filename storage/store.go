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

type query struct {
	store   Storage
	user    string
	table   string
	command sunduq.MessageType
	key     string
	value   []byte
}

func (q query) Table(table string) query {
	q.table = table
	return q
}

func (q query) User(user string) query {
	q.user = user
	return q
}

func (q query) Command(command sunduq.MessageType) query {
	q.command = command
	return q
}

func (q query) Key(key string) query {
	q.key = key

	return q
}

func (q query) Value(value []byte) query {
	q.value = value
	return q
}

func (q query) Exec() ([]byte, error) {
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
		return nil, table.Set(q.key, q.value)
	case sunduq.Del:
		if q.key == "" {
			return nil, errors.New("did not set a key for a Delete query builder")
		}
		return nil, table.Del(q.key)
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
