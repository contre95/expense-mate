package jsonstorage

import (
	"errors"
	"os"
	"sync"
)

type UserStorage struct {
	Path string
	mu   sync.Mutex
}

func NewStorage(path string) *UserStorage {
	err := CreateFileIfNotExists(path, "[]")
	if err != nil {
		panic("Can't create database" + err.Error())
	}
	return &UserStorage{
		Path: path,
	}
}

func CreateFileIfNotExists(filename, data string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			return createErr
		}
		defer file.Close()
		_, writeErr := file.Write([]byte(data))
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
