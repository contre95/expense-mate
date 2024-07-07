package jsonstorage

import (
	"encoding/json"
	"errors"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"io/fs"
	"os"
)

// All retrieves all users from the JSON storage.
func (s *UserStorage) all() ([]expense.User, error) {
	file, err := os.ReadFile(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []expense.User{}, err
		}
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	var users []expense.User
	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %w", err)
	}
	return users, nil
}

// All retrieves all users from the JSON storage.
func (s *UserStorage) All() ([]expense.User, error) {
	s.mu.Lock()
	users, err := s.all()
	defer s.mu.Unlock()
	return users, err
}

// Add adds a new user to the JSON storage.
func (s *UserStorage) Add(user expense.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.all()
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.ID == user.ID {
			return errors.New("user already exists")
		}
	}
	users = append(users, user)
	data, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("could not marshal users: %w", err)
	}
	err = os.WriteFile(s.Path, data, fs.FileMode(0644))
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}

// Delete removes a user from the JSON storage.
func (s *UserStorage) Delete(userID expense.UserID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.all()
	if err != nil {
		return err
	}
	found := false
	for i, u := range users {
		if u.ID == userID {
			users = append(users[:i], users[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not found")
	}
	data, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("could not marshal users: %w", err)
	}
	err = os.WriteFile(s.Path, data, fs.FileMode(0644))
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}
