package commands

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tinfoil-knight/tiny-redis/store"
)

var (
	ErrInvalidSyntax               = errors.New("ERR syntax error")
	ErrInvalidCommand              = errors.New("ERR unknown command")
	ErrWrongNumOfArgs              = errors.New("ERR wrong number of arguments")
	ErrValNotIntOrOutOfRange       = errors.New("ERR value is not an integer or out of range")
	ErrOffsetOutOfRange            = errors.New("ERR offset is out of range")
	ErrBitOffsetNotIntOrOutOfRange = errors.New("ERR bit offset is not an integer or out of range")
)

const NUL = "\u0000"

var EMPTY = []byte("")

func ExecuteCommand(kv *store.Store, cmdSeq []([]byte)) (res interface{}, err error) {
	s := cmdSeq
	sLen := len(s)
	cmd := strings.ToUpper(string(s[0]))
	switch cmd {
	case "PING":
		if sLen > 2 {
			return nil, ErrWrongNumOfArgs
		}
		if sLen == 2 {
			return s[1], nil
		}
		return "PONG", nil
	case "ECHO":
		if sLen != 2 {
			return nil, ErrWrongNumOfArgs
		}
		return s[1], nil
	case "GET":
		if sLen != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		if v, ok := kv.Get(key); ok {
			return v, nil
		}
		return nil, nil
	case "SET":
		if sLen < 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		v := s[2]

		if sLen > 3 {
			opts := s[2:sLen]
			nx := has(opts, "NX")
			xx := has(opts, "XX")
			get := has(opts, "GET")
			switch btoI(nx, xx, get) {
			case 2:
				if xx && get {
					r, ok := kv.Get(key)
					if ok {
						kv.Set(key, v)
					}
					return r, nil
				}
			case 1:
				if nx {
					_, ok := kv.Get(key)
					if !ok {
						kv.Set(key, v)
						return "OK", nil
					}
					return nil, nil
				} else if xx {
					_, ok := kv.Get(key)
					if ok {
						kv.Set(key, v)
						return "OK", nil
					}
					return nil, nil
				} else if get {
					r, _ := kv.Get(key)
					kv.Set(key, v)
					return r, nil
				}
			}
			return nil, ErrInvalidSyntax
		}
		kv.Set(key, v)
		return "OK", nil
	case "DEL":
		if sLen < 2 {
			return nil, ErrWrongNumOfArgs
		}
		n := 0
		for count := 1; count < sLen; count++ {
			key := s[count]
			if _, ok := kv.Get(key); ok {
				kv.Del(key)
				n++
			}
		}
		return n, nil
	case "GETDEL":
		if sLen != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		if v, ok := kv.Get(key); ok {
			kv.Del(key)
			return v, nil
		}
		return nil, nil
	case "EXISTS":
		if sLen < 2 {
			return nil, ErrWrongNumOfArgs
		}
		n := 0
		for count := 1; count < sLen; count++ {
			key := s[count]
			if _, ok := kv.Get(key); ok {
				n++
			}
		}
		return n, nil
	case "INCR":
		if sLen != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		if str, ok := kv.Get(key); ok {
			v, err := strconv.Atoi(string(str))
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v++
			kv.Set(key, []byte(strconv.Itoa(v)))
			return v, nil
		}
		kv.Set(key, []byte("1"))
		return 1, nil
	case "DECR":
		if sLen != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		if byts, ok := kv.Get(key); ok {
			v, err := strconv.Atoi(string(byts))
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v--
			kv.Set(key, []byte(strconv.Itoa(v)))
			return v, nil
		}
		kv.Set(key, []byte("-1"))
		return -1, nil
	case "INCRBY":
		if sLen != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		incr, err := strconv.Atoi(string(s[2]))
		if err != nil {
			return nil, ErrValNotIntOrOutOfRange
		}
		if byts, ok := kv.Get(key); ok {
			v, err := strconv.Atoi(string(byts))
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v += incr
			kv.Set(key, []byte(strconv.Itoa(v)))
			return v, nil
		}
		kv.Set(key, []byte(strconv.Itoa(incr)))
		return incr, nil
	case "DECRBY":
		if sLen != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		decr, err := strconv.Atoi(string(s[2]))
		if err != nil {
			return nil, ErrValNotIntOrOutOfRange
		}
		if byts, ok := kv.Get(key); ok {
			v, err := strconv.Atoi(string(byts))
			if err != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			v -= decr
			kv.Set(key, []byte(strconv.Itoa(v)))
			return v, nil
		}
		kv.Set(key, []byte(strconv.Itoa(-decr)))
		return -decr, nil
	case "APPEND":
		if sLen != 3 {
			return nil, ErrWrongNumOfArgs
		}
		// TODO: handle as bytes
		key := s[1]
		value := string(s[2])
		if v, ok := kv.Get(key); ok {
			c := string(v)
			c += value
			kv.Set(key, []byte(c))
			return len(c), nil
		}
		kv.Set(key, []byte(value))
		return len(value), nil
	case "GETBIT":
		if sLen != 3 {
			return nil, ErrWrongNumOfArgs
		}
		offset, err := strconv.Atoi(string(s[2]))
		if err != nil || offset < 0 {
			return nil, ErrBitOffsetNotIntOrOutOfRange
		}
		key := s[1]
		v, ok := kv.Get(key)
		if !ok {
			return 0, nil
		}
		if offset >= len(v) || s[1][offset] == 0 {
			return 0, nil
		}
		return 1, nil
	case "SETBIT":
	case "SAVE":
		kv.Save()
		return "OK", nil
	case "STRLEN":
		if sLen != 2 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		v, _ := kv.Get(key)
		return len(v), nil
	case "GETRANGE":
		if sLen != 4 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		if v, ok := kv.Get(key); ok {
			l := len(v)
			start, err1 := strconv.Atoi(string(s[2]))
			end, err2 := strconv.Atoi(string(s[3]))
			if err1 != nil || err2 != nil {
				return nil, ErrValNotIntOrOutOfRange
			}
			if start >= l {
				return EMPTY, nil
			}
			if end >= l {
				end = l - 1
			}
			start = (start%l + l) % l
			end = (end%l + l) % l
			if start > end {
				return EMPTY, nil
			}
			// GETRANGE is inclusive for both offsets
			end++
			return []byte(v[start:end]), nil
		}
		return EMPTY, nil
	case "SETRANGE":
		if sLen != 4 {
			return nil, ErrWrongNumOfArgs
		}
		// TODO: handle as bytes
		key := s[1]
		offset, err := strconv.Atoi(string(s[2]))
		if err != nil {
			return nil, ErrValNotIntOrOutOfRange
		}
		value := string(s[3])
		v, _ := kv.Get(key)
		if offset < 0 {
			return nil, ErrOffsetOutOfRange
		}
		c := string(v)
		if offset >= len(v) {
			d := offset - len(v)
			for i := 0; i < d; i++ {
				c += fmt.Sprintf("%v", NUL)
			}
			c += value
		} else {
			if len(value) < len(v) {
				c = c[:offset] + value + c[len(value)+offset:]
			} else {
				c = c[:offset] + value
			}
		}
		kv.Set(key, []byte(c))
		return len(c), nil
	case "MGET":
		if sLen < 2 {
			return nil, ErrWrongNumOfArgs
		}
		keys := s[1:]
		r := make([][]byte, len(keys))
		for i, key := range keys {
			v, ok := kv.Get(key)
			if ok {
				r[i] = v
			} else {
				r[i] = []byte(nil)
			}
		}
		return r, nil
	case "MSET":
		if sLen < 3 || (sLen&1 == 0) {
			return nil, ErrWrongNumOfArgs
		}
		pairs := s[1:]
		for i := 0; i < len(pairs)-1; i += 2 {
			kv.Set(pairs[i], pairs[i+1])
		}
		return "OK", nil
	case "MSETNX":
		if sLen < 3 || (sLen&1 == 0) {
			return nil, ErrWrongNumOfArgs
		}
		// TODO: make the operation atomic
		pairs := s[1:]
		n := 0
		for i := 0; i < len(pairs)-1; i += 2 {
			_, ok := kv.Get(pairs[i])
			if ok {
				n++
			}
		}
		if n > 0 {
			return 0, nil
		}
		for i := 0; i < len(pairs)-1; i += 2 {
			kv.Set(pairs[i], pairs[i+1])
		}
		return 1, nil
	case "COPY":
		if sLen < 3 || sLen > 4 {
			return nil, ErrWrongNumOfArgs
		}
		src := s[1]
		v, ok := kv.Get(src)
		if !ok {
			return 0, nil
		}
		dest := s[2]
		_, ok = kv.Get(dest)
		if ok {
			if sLen == 4 {
				if !bytes.Equal(s[3], []byte("REPLACE")) {
					return nil, ErrInvalidSyntax
				}
			} else {
				return 0, nil
			}
		}
		kv.Set(dest, v)
		return 1, nil
	case "RENAME":
	case "FLUSHDB":
	}
	return nil, ErrInvalidCommand
}

func has(arr [][]byte, value string) bool {
	for _, v := range arr {
		if strings.ToUpper(string(v)) == value {
			return true
		}
	}
	return false
}

func btoI(b ...bool) int {
	n := 0
	for _, v := range b {
		if v {
			n++
		}
	}
	return n
}
