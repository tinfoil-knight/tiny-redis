// Package resp implements functions for decoding and encoding RESP 2.
package resp

import (
	"errors"
	"strconv"
)

var (
	ErrInvalidDataType = errors.New("invalid datatype")
)

const (
	SIMPLE_STRING = '+'
	ERROR         = '-'
	INTEGER       = ':'
	BULK_STRING   = '$'
	ARRAY         = '*'
)

func Decode(input []byte) (decodedValue interface{}, read int) {
	switch first_byte := input[0]; first_byte {
	case SIMPLE_STRING:
		return handleSimpleString(input[1:])
	case ERROR:
		return handleError(input[1:])
	case INTEGER:
		return handleInteger(input[1:])
	case BULK_STRING:
		return handleBulkString(input[1:])
	case ARRAY:
		return handleArray(input[1:])
	}
	return ErrInvalidDataType, 1
}

func readUntilCRLF(bytes []byte) (string, int) {
	str := ""
	var c byte
	read := 0
	for i := 0; i < len(bytes); i++ {
		c = bytes[i]
		read++
		if c == '\n' {
			break
		}
		if c != '\r' {
			str += string(c)
		}
	}
	return str, read
}

func handleSimpleString(in []byte) (string, int) {
	v, read := readUntilCRLF(in)
	return v, read
}

func handleError(in []byte) (error, int) {
	v, read := readUntilCRLF(in)
	return errors.New(v), read
}

func handleInteger(in []byte) (int, int) {
	str, read := readUntilCRLF(in)
	v, _ := strconv.Atoi(str)
	return v, read
}

func handleBulkString(in []byte) (interface{}, int) {
	length, read := readUntilCRLF(in)
	size, _ := strconv.Atoi(length)
	switch size {
	case 0:
		return []byte(""), read
	case -1:
		return nil, read
	default:
		val, r := readUntilCRLF(in[read:])
		return []byte(val), read + r
	}
}

func handleArray(in []byte) (interface{}, int) {
	length, read := readUntilCRLF(in)
	size, _ := strconv.Atoi((length))
	switch size {
	case 0:
		empty := make([]interface{}, 0)
		return empty, read
	case -1:
		return nil, read
	default:
		in = in[read:]
		totalRead := read
		items := make([]interface{}, size)
		for counter := 0; counter < size; counter++ {
			item, r := Decode([]byte(in))
			// First byte is skipped in Decode
			totalRead += r + 1
			in = in[r+1:]
			items[counter] = item
		}
		return items, totalRead
	}
}
