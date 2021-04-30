package resp

import (
	"fmt"
	"reflect"
)

const NIL = "_\r\n"

func Encode(input interface{}) string {
	switch input.(type) {
	case int:
		return fmt.Sprintf(":%d\r\n", input)
	case string:
		return fmt.Sprintf("+%s\r\n", input)
	case []byte:
		len := len(reflect.ValueOf(input).Bytes())
		return fmt.Sprintf("$%v\r\n%s\r\n", len, input)
	case error:
		return fmt.Sprintf("-%s\r\n", input)
	case [][]byte:
		s := reflect.ValueOf(input)
		v := fmt.Sprintf("*%v\r\n", s.Len())
		for i := 0; i < s.Len(); i++ {
			b := (s.Index(i).Interface()).([]byte)
			if len(b) > 0 {
				v += Encode(b)
			} else {
				v += Encode(nil)
			}
		}
		return v
	case nil:
		return NIL
	}
	panic(ErrInvalidSyntax)
}
