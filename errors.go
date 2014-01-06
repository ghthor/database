package database

import (
	"errors"
	"fmt"
)

type Err struct {
	Err error
}

func (e Err) Error() string {
	return fmt.Sprintf("database error:%s", e.Err.Error())
}

var (
	ErrUnimplemented     = errors.New("unimplemented")
	ErrInvalidAction     = errors.New("attempted to execute with an invalid action")
	ErrResourceForbidden = errors.New("resource forbidden")
)
