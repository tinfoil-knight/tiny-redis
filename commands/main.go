package commands

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrInvalidCommand        = errors.New("ERR unknown command")
	ErrWrongNumOfArgs        = errors.New("ERR wrong number of arguments")
	ErrValNotIntOrOutOfRange = errors.New("ERR value is not an integer or out of range")
)
var kv = make(map[string]string)

func ExecuteCommand(arr interface{}) (interface{}, error) {
	s := reflect.ValueOf(arr)
	cmd := fmt.Sprintf("%s", s.Index(0))
	switch cmd {
	case "PING":
		if s.Len() > 2 {
			return nil, ErrWrongNumOfArgs
		}
		if s.Len() == 2 {
			arg := fmt.Sprintf("%s", s.Index(1))
			return []byte(arg), nil
		}
		return "PONG", nil
	case "ECHO":
		if s.Len() != 2 {
			return nil, ErrWrongNumOfArgs
		}
		arg := fmt.Sprintf("%s", s.Index(1))
		return []byte(arg), nil
	case "GET":
		if s.Len() != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		if v, ok := kv[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	case "SET":
		if s.Len() != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		kv[key] = fmt.Sprintf("%s", s.Index(2))
		return "OK", nil
	case "DEL":
		if s.Len() < 2 {
			return nil, ErrWrongNumOfArgs
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
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		if v, ok := kv[key]; ok {
			delete(kv, key)
			return []byte(v), nil
		}
		return nil, nil
	case "EXISTS":
		if s.Len() < 2 {
			return nil, ErrWrongNumOfArgs
		}
		n := 0
		for count := 1; count < s.Len(); count++ {
			if _, ok := kv[fmt.Sprintf("%s", s.Index(count))]; ok {
				n++
			}
		}
		return n, nil
	case "INCR":
		if s.Len() != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		if str, ok := kv[key]; ok {
			v, err := strconv.Atoi(str)
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v++
			kv[key] = strconv.Itoa(v)
			return v, nil
		}
		kv[key] = "1"
		return 1, nil
	case "DECR":
		if s.Len() != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		if str, ok := kv[key]; ok {
			v, err := strconv.Atoi(str)
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v--
			kv[key] = strconv.Itoa(v)
			return v, nil
		}
		kv[key] = "-1"
		return -1, nil
	case "INCRBY":
	case "DECRBY":
	case "APPEND":
	case "GETBIT":
	case "SETBIT":
	case "QUIT":
	case "SAVE":
	case "STRLEN":
	case "GETRANGE":
	case "SETRANGE":
	case "SETNX":
	case "MGET":
	case "MSET":
	case "MSETNX":
	case "RENAME":

	default:
		return nil, ErrInvalidCommand
	}
	return nil, errors.New("commands.ExecuteCommand: unknown error")
}
