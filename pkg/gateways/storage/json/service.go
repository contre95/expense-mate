package jsonStorage

import "os/user"

type JSONStorage struct {
}

func NewStorage() *JSONStorage {
	return &JSONStorage{}
}

func (js *JSONStorage) Add(u user.User) error {
	panic("not implemented") // TODO: Implement
}

func (js *JSONStorage) Get(uname string) (*user.User, error) {
	panic("not implemented") // TODO: Implement
}
