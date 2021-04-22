package resp

import (
	"fmt"
	"reflect"
)

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
	case []interface{}:
		s := reflect.ValueOf(input)
		v := fmt.Sprintf("*%v\r\n", s.Len())
		for i := 0; i < s.Len(); i++ {
			v += Encode(s.Index(i).Interface())
		}
		return v
	case nil:
		return "_\r\n"
	default:
		panic(ErrInvalidDataType)
	}
}
