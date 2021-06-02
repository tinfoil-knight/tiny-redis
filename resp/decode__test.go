package resp

import (
	"errors"
	"reflect"
	"testing"
)

func Test__SimpleString(t *testing.T) {
	got, _ := Decode([]byte("+OK\r\n"))
	want := "OK"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Test__Integer(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{":1000\r\n", 1000},
		{":-1\r\n", -1},
		{":0\r\n", 0},
	}
	for _, tt := range tests {
		got, _ := Decode([]byte(tt.input))
		if got != tt.expected {
			t.Errorf("Decode(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test__BulkString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"$6\r\nfoobar\r\n", []byte("foobar")},
		{"$0\r\n\r\n", []byte("")},
	}

	for _, tt := range tests {
		got, _ := Decode([]byte(tt.input))
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Decode(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test__Array(t *testing.T) {
	emptySlice := make([]interface{}, 0)
	sliceOfStrings := []interface{}{"foo", "bar"}
	sliceOfIntegers := []interface{}{1, 2, 3}
	sliceWithMixedTypes := []interface{}{1, 2, 3, 4, []byte("foobar")}
	nestedSlice := []interface{}{sliceOfIntegers, sliceOfStrings}
	sliceWithNull := []interface{}{[]byte("foo"), nil, []byte("bar")}
	sliceWithEmptyString := []interface{}{[]byte("foo"), []byte(""), []byte("bar")}

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"*0\r\n", emptySlice},
		{"*2\r\n+foo\r\n+bar\r\n", sliceOfStrings},
		{"*3\r\n:1\r\n:2\r\n:3\r\n", sliceOfIntegers},
		{"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n", sliceWithMixedTypes},
		{"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+foo\r\n+bar\r\n", nestedSlice},
		{"*3\r\n$3\r\nfoo\r\n_\r\n$3\r\nbar\r\n", sliceWithNull},
		{"*3\r\n$3\r\nfoo\r\n$0\r\n\r\n$3\r\nbar\r\n", sliceWithEmptyString},
	}
	for _, tt := range tests {
		got, _ := Decode([]byte(tt.input))
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Decode(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test__Error(t *testing.T) {
	got, _ := Decode([]byte("-Error message\r\n"))
	want := errors.New("Error message")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
