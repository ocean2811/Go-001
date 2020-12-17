package data

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserInfo struct {
	ID   string
	Name string
}

type DataOperator interface {
	GetUserInfo(id string) (*UserInfo, error)
}
