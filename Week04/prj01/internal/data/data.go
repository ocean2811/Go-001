package data

import "github.com/pkg/errors"

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
