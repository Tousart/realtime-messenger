package hashpassword

import (
	"testing"
)

func TestBCryptPasswordHasher(t *testing.T) {
	ph := NewBCryptPasswordHasher()
	password := "123456789"

	hashPassword, err := ph.Hash(password)
	if err != nil {
		t.Errorf("hash password error: %v", err)
	}

	expected := true
	actual := ph.Compare(hashPassword, password)

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
