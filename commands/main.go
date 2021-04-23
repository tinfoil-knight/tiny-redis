package commands

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidCommand         = errors.New("ERR unknown command")
	ErrWrongNumberOfArguments = errors.New("ERR wrong number of arguments")
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
	case "DEL":
		if s.Len() < 2 {
			return nil, ErrWrongNumberOfArguments
		}
		n := 0
		for count := 1; count < s.Len(); count++ {
			key := fmt.Sprintf("%s", s.Index(count))
			if _, ok := kv[key]; ok {
				delete(kv, key)
				n++
			}
		}
		return n, nil
	case "GETDEL":
		if s.Len() != 2 {
			return nil, ErrWrongNumberOfArguments
		}
		key := fmt.Sprintf("%s", s.Index(1))
		if v, ok := kv[key]; ok {
			delete(kv, key)
			return []byte(v), nil
		}
		return nil, nil
	case "EXISTS":
		n := 0
		for count := 1; count < s.Len(); count++ {
			if _, ok := kv[fmt.Sprintf("%s", s.Index(count))]; ok {
				n++
			}
		}
		return n, nil
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
