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
var kv = make(map[string]string)

func ExecuteCommand(arr interface{}) (interface{}, error) {
	s := reflect.ValueOf(arr)
	cmd := fmt.Sprintf("%s", s.Index(0))
	switch cmd {
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
		if s.Len() != 2 {
			return nil, ErrWrongNumberOfArguments
		}
		if v, ok := kv[fmt.Sprintf("%s", s.Index(1))]; ok {
			return []byte(v), nil
		}
		return nil, nil
	case "SET":
		if s.Len() != 3 {
			return nil, ErrWrongNumberOfArguments
		}
		kv[fmt.Sprintf("%s", s.Index(1))] = fmt.Sprintf("%s", s.Index(2))
		return "OK", nil
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
