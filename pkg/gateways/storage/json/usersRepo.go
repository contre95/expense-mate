package json

import (
	"encoding/json"
	"expenses-app/pkg/domain/user"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
)

func (js *JSONStorage) getAllUsers() ([]jsonUser, error) {
	var users []byte
	var err error
	users, err = ioutil.ReadFile(js.Path)
	if err != nil {
		return nil, err
	}
	var allUsers []jsonUser
	err = json.Unmarshal(users, &allUsers)
	if err != nil {
		return nil, err
	}
	return allUsers, nil
}

func (js *JSONStorage) Add(u user.User) error {
	var allUsers []jsonUser
	var oldFile []byte
	var newFile []byte
	var err error
	oldFile, err = ioutil.ReadFile(js.Path)
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(oldFile, &allUsers)
	if err != nil {
		log.Println(err)
		return err
	}
	newJsonUser := jsonUser{
		Username: u.Username,
		Alias:    u.Alias,
		UUID:     u.ID.String(),
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

func (js *JSONStorage) Exists(uname string) (bool, error) {
	user, err := js.Get(uname)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (js *JSONStorage) Get(uname string) (*user.User, error) {
	jsonUsers, err := js.getAllUsers()
	if err != nil {
		return nil, err
	}
	for _, jsonUser := range jsonUsers {
		if jsonUser.Username == uname {
			uuid, err := uuid.Parse(jsonUser.UUID)
			if err != nil {
				return nil, err
			}
			return &user.User{
				ID:       uuid,
				Username: jsonUser.Username,
				Password: jsonUser.Password,
				Alias:    jsonUser.Alias,
			}, nil
		}
	}
	return nil, nil
}
