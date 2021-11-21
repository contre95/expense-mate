package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}
func (ph *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (ph *PasswordHasher) CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
