package db

import (
	"errors"
	"io/ioutil"
	"os"
)

type FileDB struct {
	root string
}

func NewFileDB(name string) (*FileDB, error) {
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(name, 0777)
			if err != nil {
				return nil, errors.New("Cannot create DB: " + name + " with error" + err.Error())
			}
		}
	} else {
		if !info.IsDir() {
			return nil, errors.New("File " + name + " Already exist! And it is not a directory")
		}
	}

	db := &FileDB{root: name}

	return db, nil
}

func (fdb *FileDB) Save(key string, value []byte) error {
	file, err := os.Create(fdb.root + "/" + key)
	defer file.Close()
	if err != nil {
		return errors.New("Cannot save: " + key + " with error: " + err.Error())
	}

	if _, err := file.Write(value); err != nil {
		return errors.New("Cannot save: " + key + " with error: " + err.Error())
	}

	return nil
}

func (fdb *FileDB) Get(key string) ([]byte, error) {
	value, err := ioutil.ReadFile(fdb.root + "/" + key)
	if err != nil {
		return nil, errors.New("Cannot get: " + key + " with error: " + err.Error())
	}

	return value, nil
}

func (fdb *FileDB) Destroy() error {
	err := os.RemoveAll(fdb.root)
	if err != nil {
		return errors.New("Cannot Destroy DB: " + fdb.root + " with error: " + err.Error())
	}

	return nil
}
