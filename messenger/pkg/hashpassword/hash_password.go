package pkg

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type BCryptPasswordHasher struct{}

func NewBCryptPasswordHasher() *BCryptPasswordHasher {
	return &BCryptPasswordHasher{}
}

func (ph *BCryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("pkg: Hash: %s", err.Error())
	}
	return string(bytes), nil
}

func (ph *BCryptPasswordHasher) Compare(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
