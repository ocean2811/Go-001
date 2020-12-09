package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

//1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

func main() {
	Biz()
}

// User info
type User struct {
	ID   uint64
	Info struct{}
}

// Biz handle rpc request
func Biz() {
	id := uint64(9008000000048942)
	user, err := Dao(id)
	if err != nil {
		fmt.Printf("Error: Dao has error.userID=%d,err=%+v\n", id, err)
		return
	}

	fmt.Printf("Info: login complete.userID=%d,info=%+v\n", id, user)
}

// Dao query user info
func Dao(id uint64) (*User, error) {
	//access DB...
	err := sql.ErrNoRows
	return nil, errors.Errorf("query user info has error.userID=%d,err=%v", id, err)
}
