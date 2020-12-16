package service

import (
	"errors"
	"net"
)

var (
	ErrServerUnavailable = errors.New("server unavailable")
	ErrUserInvalid       = errors.New("user invalid")
)

type Server interface {
	Serve(net.Listener) error
	Stop()
}
