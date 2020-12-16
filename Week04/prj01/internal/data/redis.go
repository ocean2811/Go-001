package data

import (
	"time"

	"github.com/pkg/errors"
)

type redisUserInfoClient struct {
	addr string
}

func NewRedis(addr string) DataOperator {
	return &redisUserInfoClient{addr: addr}
}

func (redis *redisUserInfoClient) GetUserInfo(id string) (*UserInfo, error) {
	if time.Now().Unix()%2 == 0 {
		//data access has error
		return nil, errors.WithMessagef(ErrUserNotFound, "user_id=%s", id)
	}

	//data access success
	return &UserInfo{ID: id, Name: "Fake Man"}, nil
}
