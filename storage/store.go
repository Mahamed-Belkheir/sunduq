package storage

import (
	"errors"
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

//Add adds data to a specific table
func (s Storage) Add(req sunduq.Request, tableName, key string, value []byte) error {
	defer s.mut.RUnlock()
	s.mut.RLock()
	table, ok := s.tables[tableName]
	if !ok {
		return tableNotFound()
	}
	return table.Add(req, key, value)
}

//AddTable creates and adds a new table to the storage struct
func (s Storage) AddTable(req sunduq.Request, tableName string, isSystemTable bool) error {
	defer s.mut.Unlock()
	s.mut.Lock()
	_, ok := s.tables[tableName]
	if ok {
		return tableExists()
	}
	s.tables[tableName] = NewTable(tableName, req.User, isSystemTable)
	return nil
}
