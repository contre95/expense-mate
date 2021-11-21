package json

import (
	"encoding/json"
	"expenses-app/pkg/domain/user"
	"io/ioutil"
)

func (js *JSONStorage) Add(u user.User) error {
	var allUsers []jsonUser
	var oldFile []byte
	var newFile []byte
	var err error
	oldFile, err = ioutil.ReadFile(js.Path + "/users.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(oldFile, &allUsers)
	if err != nil {
		return err
	}
	newJsonUser := jsonUser{
		Username: u.Password,
		Alias:    u.Alias,
		Password: u.Password,
	}
	allUsers = append(allUsers, newJsonUser)

	newFile, err = json.MarshalIndent(allUsers, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(js.Path, newFile, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (js *JSONStorage) Get(uname string) (*user.User, error) {
	panic("not implemented") // TODO: Implement
}
