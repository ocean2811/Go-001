package biz

import (
	"go001/internal/data"

	"github.com/pkg/errors"
)

var (
	ErrUninitialized = errors.New("biz dataOP uninitialized")
)

type BizHandler struct {
	dataOP data.DataOperator
}

func NewBiz(dataOp data.DataOperator) *BizHandler {
	return &BizHandler{dataOP: dataOp}
}

func (b *BizHandler) GenHelloMsg(id string) (string, error) {
	if b.dataOP == nil {
		return "", errors.WithMessagef(ErrUninitialized, "id=%s", id)
	}

	info, err := b.dataOP.GetUserInfo(id)
	if err != nil {
		return "", err
	}

	return info.Name + " Hello!!!", nil
}
