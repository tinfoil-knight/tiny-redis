package resp

import (
	"errors"
	"testing"
)

func Test__StringEn(t *testing.T) {
	got := Encode("foobar")
	want := "+foobar\r\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Test__BulkStringEn(t *testing.T) {
	got := Encode([]byte("foobar"))
	want := "$6\r\nfoobar\r\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Test__IntegerEn(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1000, ":1000\r\n"},
		{-1, ":-1\r\n"},
		{0, ":0\r\n"},
	}
	for _, tt := range tests {
		got := Encode(tt.input)
		if got != tt.expected {
			t.Errorf("Encode(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test__ErrorEn(t *testing.T) {
	got := Encode(errors.New("Error message"))
	want := "-Error message\r\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Test__ArrayEn(t *testing.T) {
	emptySlice := make([]interface{}, 0)
	sliceOfStrings := []interface{}{"foo", "bar"}
	sliceOfIntegers := []interface{}{1, 2, 3}
	sliceWithMixedTypes := []interface{}{1, 2, 3, 4, []byte("foobar")}
	nestedSlice := []interface{}{sliceOfIntegers, sliceOfStrings}
	sliceWithNull := []interface{}{[]byte("foo"), nil, []byte("bar")}

	tests := []struct {
		input    interface{}
		expected string
	}{
		{emptySlice, "*0\r\n"},
		{sliceOfStrings, "*2\r\n+foo\r\n+bar\r\n"},
		{sliceOfIntegers, "*3\r\n:1\r\n:2\r\n:3\r\n"},
		{sliceWithMixedTypes, "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n"},
		{nestedSlice, "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+foo\r\n+bar\r\n"},
		{sliceWithNull, "*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n"},
	}
	for _, tt := range tests {
		got := Encode(tt.input)
		if got != tt.expected {
			t.Errorf("Encode(%v): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test_NilEn(t *testing.T) {
	got := Encode(nil)
	want := "$-1\r\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
