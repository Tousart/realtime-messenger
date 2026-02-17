package domain

import "strings"

type User struct {
	UserID   int
	UserName string
	Password string
}

type UserField func(user *User) error

func NewUser(fields ...UserField) (*User, error) {
	user := User{}
	for _, field := range fields {
		err := field(&user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func WithUserName(userName string) UserField {
	return func(user *User) error {
		name := strings.TrimSpace(userName)
		if len(name) < 2 || len(name) > 30 {
			return ErrBadUserName
		}
		user.UserName = name
		return nil
	}
}

func WithPassword(password string) UserField {
	return func(user *User) error {
		pswrd := strings.TrimSpace(password)
		if len(pswrd) < 8 || len(pswrd) > 72 {
			return ErrBadPassword
		}
		user.Password = pswrd
		return nil
	}
}
