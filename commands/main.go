package commands

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tinfoil-knight/tiny-redis/store"
)

var (
	ErrInvalidCommand        = errors.New("ERR unknown command")
	ErrWrongNumOfArgs        = errors.New("ERR wrong number of arguments")
	ErrValNotIntOrOutOfRange = errors.New("ERR value is not an integer or out of range")
	ErrOffsetOutOfRange      = errors.New("ERR offset is out of range")
)

const NUL = "\u0000"

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
		if sLen != 3 {
			return nil, ErrWrongNumOfArgs
		}
		key := s[1]
		v := s[2]
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
	case "SETBIT":
	case "SAVE":
		f, err := os.OpenFile("dump.trdb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		b := new(bytes.Buffer)
		if err = gob.NewEncoder(b).Encode(kv.GetUnderlying()); err != nil {
			panic(err)
		}
		if _, err = io.Copy(f, b); err != nil {
			panic(err)
		}
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
	case "SETNX":
	case "MGET":
	case "MSET":
	case "MSETNX":
	case "RENAME":
	case "FLUSHDB":
	}
	// TODO(fix): This error isn't sent back
	return nil, ErrInvalidCommand
}
