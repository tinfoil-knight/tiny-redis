package commands

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
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
			key := fmt.Sprintf("%s", s.Index(count))
			if _, ok := kv[key]; ok {
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
		if s.Len() != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		incr, err := strconv.Atoi(fmt.Sprintf("%s", s.Index(2)))
		if err != nil {
			return nil, ErrValNotIntOrOutOfRange
		}
		if str, ok := kv[key]; ok {
			v, err := strconv.Atoi(str)
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v += incr
			kv[key] = strconv.Itoa(v)
			return v, nil
		}
		kv[key] = strconv.Itoa(incr)
		return incr, nil
	case "DECRBY":
		if s.Len() != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		decr, err := strconv.Atoi(fmt.Sprintf("%s", s.Index(2)))
		if err != nil {
			return nil, ErrValNotIntOrOutOfRange
		}
		if str, ok := kv[key]; ok {
			v, err := strconv.Atoi(str)
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v -= decr
			kv[key] = strconv.Itoa(v)
			return v, nil
		}
		kv[key] = strconv.Itoa(-decr)
		return -decr, nil
	case "APPEND":
		if s.Len() != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		value := fmt.Sprintf("%s", s.Index(2))
		if v, ok := kv[key]; ok {
			v += value
			kv[key] = v
			return len(v), nil
		}
		kv[key] = value
		return len(value), nil
	case "GETBIT":
	case "SETBIT":
	case "SAVE":
		f, err := os.OpenFile("dump.tdb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		b := new(bytes.Buffer)
		if err = gob.NewEncoder(b).Encode(kv); err != nil {
			panic(err)
		}
		if _, err = io.Copy(f, b); err != nil {
			panic(err)
		}
		return "OK", nil
	case "STRLEN":
	case "GETRANGE":
		if s.Len() != 4 {
			return nil, ErrWrongNumOfArgs
		}
		key := fmt.Sprintf("%s", s.Index(1))
		if v, ok := kv[key]; ok {
			l := len(v)
			start, _ := strconv.Atoi(fmt.Sprintf("%s", s.Index(2)))
			end, _ := strconv.Atoi(fmt.Sprintf("%s", s.Index(3)))
			if start >= l {
				return []byte(""), nil
			}
			if end >= l {
				end = l - 1
			}
			start = (start%l + l) % l
			end = (end%l + l) % l
			if start > end {
				return []byte(""), nil
			}
			// GETRANGE is inclusive for both offsets
			end++
			return []byte(v[start:end]), nil
		}
		return []byte(""), nil
	case "SETRANGE":
	case "SETNX":
	case "MGET":
	case "MSET":
	case "MSETNX":
	case "RENAME":
	case "FLUSHDB":
	default:
		return nil, ErrInvalidCommand
	}
	return nil, errors.New("commands.ExecuteCommand: unknown error")
}

// func load() {
// 	f, err := os.Open("dump.tdb")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	if err = gob.NewDecoder(f).Decode(&kv); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%#v\n", kv)
// }
