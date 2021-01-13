package storage

import (
	"io/ioutil"
)

//PersistenceAPI struct holds all the methods to handle storing and retrieving buckets from disk
type PersistenceAPI struct {
	filePath string
}

//Save saves the data into disk
func (p *PersistenceAPI) Save(data []byte) error {
	return ioutil.WriteFile(p.filePath, data, 0644)
}

//Load loads the table data from disk
func (p *PersistenceAPI) Load() ([]byte, error) {
	return ioutil.ReadFile(p.filePath)
}
