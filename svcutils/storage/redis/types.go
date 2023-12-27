package redis

import (
	"errors"
	"fmt"
)

type any = interface{}
type nullObj struct{}

func NilObj() *nullObj {
	return &nullObj{}
}

type dbloc int

func (d dbloc) ToInt() int {
	return int(d)
}

const (
	DefaultDB dbloc = iota
)

var ErrNotFound = errors.New("key not found")
var ErrNilDataStorage = errors.New("data cannot be nil")
var ErrNonDataPtr = errors.New("must pass a pointer, not a value, to StructScan destination")

type Keys string

func (k Keys) ToString() string {
	return string(k)
}

func Session(uid string) Keys {
	return Keys(fmt.Sprintf("session-%v", uid))
}
