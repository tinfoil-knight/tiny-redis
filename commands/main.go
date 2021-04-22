package commands

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidCommand         = errors.New("no such command found")
	ErrWrongNumberOfArguments = errors.New("wrong number of arguments")
)

func ExecuteCommand(arr interface{}) (interface{}, error) {
	s := reflect.ValueOf(arr)

	switch fmt.Sprintf("%s", s.Index(0)) {
	case "PING":
		if s.Len() > 2 {
			return nil, ErrWrongNumberOfArguments
		} else {
			if s.Len() == 2 {
				return []byte(fmt.Sprintf("%s", s.Index(1))), nil
			}
			return "PONG", nil
		}
	case "GET":
	case "SET":
	case "GETSET":
	case "DEL":
	case "GETDEL":
	case "EXISTS":
	case "INCR":
	case "INCRBY":
	case "DECR":
	case "DECRBY":
	case "QUIT":
	case "SAVE":
	case "STRLEN":
	case "RENAME":

	default:
		return nil, ErrInvalidCommand
	}
	return nil, errors.New("commands.ExecuteCommand: unknown error")
}
