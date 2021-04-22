// Package resp implements functions for decoding and encoding RESP 3.
package resp

import (
	"errors"
	"math/big"
	"strconv"
)

var (
	ErrInvalidDataType = errors.New("invalid datatype")
	ErrInvalidInput    = errors.New("invalid input")
)

const (
	SIMPLE_STRING   = '+'
	ERROR           = '-'
	INTEGER         = ':'
	DOUBLE          = ','
	BIGINT          = '('
	BOOLEAN         = '#'
	BULK_STRING     = '$'
	BULK_ERROR      = '!'
	VERBATIM_STRING = '='
	ARRAY           = '*'
	SET             = '~'
	NULL            = '_'
)

func Decode(input []byte) (decodedValue interface{}, read int) {
	switch first_byte := input[0]; first_byte {
	case SIMPLE_STRING:
		return handleSimpleString(input[1:])
	case ERROR:
		return handleError(input[1:])
	case INTEGER:
		return handleInteger(input[1:])
	case DOUBLE:
		return handleDouble(input[1:])
	case BIGINT:
		return handleBigInt(input[1:])
	case BOOLEAN:
		return handleBoolean(input[1:])
	case BULK_STRING, VERBATIM_STRING:
		return handleBulkString(input[1:])
	case BULK_ERROR:
		return handleBulkError(input[1:])
	case ARRAY:
		return handleArray(input[1:])
	case SET:
		return handleSet(input[1:])
	case NULL:
		return nil, 3
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
	// TODO: Handle -> First word must be upper-case
	v, read := readUntilCRLF(in)
	return errors.New(v), read
}

func handleInteger(in []byte) (int, int) {
	str, read := readUntilCRLF(in)
	v, _ := strconv.Atoi(str)
	return v, read
}

func handleDouble(in []byte) (float64, int) {
	str, read := readUntilCRLF(in)
	// TODO: Handle inf
	v, _ := strconv.ParseFloat(str, 64)
	return v, read
}

func handleBigInt(in []byte) (*big.Int, int) {
	str, read := readUntilCRLF(in)
	v, _ := new(big.Int).SetString(str, 10)
	return v, read
}

func handleBoolean(in []byte) (bool, int) {
	str, read := readUntilCRLF(in)
	switch str {
	case "f":
		return false, read
	case "t":
		return true, read
	default:
		panic(ErrInvalidInput)
	}
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

func handleBulkError(in []byte) (error, int) {
	v, read := handleBulkString(in)
	return errors.New(v.(string)), read
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

func handleSet(in []byte) (interface{}, int) {
	length, read := readUntilCRLF(in)
	size, _ := strconv.Atoi((length))

	set := make(map[interface{}]bool)

	switch size {
	case 0:
		return set, read
	case -1:
		return nil, read
	default:
		in = in[read:]
		totalRead := read
		for counter := 0; counter < size; counter++ {
			item, r := Decode([]byte(in))
			// First byte is skipped in Decode
			totalRead += r + 1
			in = in[r+1:]
			set[item] = true
		}
		return set, totalRead
	}
}
