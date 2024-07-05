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
	err := createFileIfNotExists(path)
	if err != nil {
		panic("Can't create database" + err.Error())
	}
	return &UserStorage{
		Path: path,
	}
}

func createFileIfNotExists(filename string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			return createErr
		}
		defer file.Close()
		_, writeErr := file.Write([]byte("[]"))
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
