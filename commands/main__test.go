package commands

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/tinfoil-knight/tiny-redis/store"
)

func b(x string) []byte {
	return []byte(x)
}

func bA(x []string) [][]byte {
	arr := make([][]byte, len(x))
	for i := range arr {
		arr[i] = []byte(x[i])
	}
	return arr
}

func Test__PING_ECHO(t *testing.T) {
	kv := store.New()
	tests := []struct {
		input    []([]byte)
		expected interface{}
	}{
		{bA([]string{"PING"}), "PONG"},
		{bA([]string{"PING", "hey"}), "hey"},
		{bA([]string{"ECHO", "hey"}), "hey"},
	}
	for _, tt := range tests {
		v, _ := ExecuteCommand(kv, tt.input)
		got := fmt.Sprintf("%s", v)
		if got != tt.expected {
			t.Errorf("ExecuteCommand(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test__SET(t *testing.T) {
	kv := store.New()
	s := "Hello World"
	s1 := "foobar"
	for i := 1; i < 5; i++ {
		kv.Set(b(fmt.Sprintf("preset%d", i)), b(s))
	}
	type setTest struct {
		name     string
		key      string
		value    string
		options  []string
		expected string
	}
	withValueChk := []setTest{
		{"No Options Simple Key", "notset1", s, []string{}, s},
		{"No Options Empty Key", "", s, []string{}, s},
		{"New Key with NX", "notset2", s, []string{"NX"}, s},
		{"Preset Key with NX", "preset1", s1, []string{"NX"}, s},
		{"New Key with XX", "notset3", s, []string{"XX"}, ""},
		{"Preset Key with XX", "preset2", s1, []string{"XX"}, s1},
	}
	withReturnChk := []setTest{
		{"New Key with GET", "notset4", s1, []string{"GET"}, ""},
		{"Preset Key with GET", "preset3", s1, []string{"GET"}, s},
		{"New Key with XX,GET", "notset5", s, []string{"XX", "GET"}, ""},
		{"Preset Key with XX,GET", "preset4", s1, []string{"XX", "GET"}, s},
		// skip: all other cases with 2 keys which lead to syntax error
	}
	for _, tt := range withValueChk {
		in := append([]string{"SET", tt.key, tt.value}, tt.options...)
		ExecuteCommand(kv, bA(in))
		v, _ := kv.Get(b(tt.key))
		got := fmt.Sprintf("%s", v)
		if got != tt.expected {
			t.Errorf("ExecuteCommand(%q)(%q): got %q want %q", tt.name, in, got, tt.expected)
		}
	}
	for _, tt := range withReturnChk {
		in := append([]string{"SET", tt.key, tt.value}, tt.options...)
		got, _ := ExecuteCommand(kv, bA(in))
		if fmt.Sprintf("%s", got) != tt.expected {
			t.Errorf("ExecuteCommand(%q)(%s): got %q want %q", tt.name, in, got, tt.expected)
		}
	}
}

func Test__GETRANGE(t *testing.T) {
	kv := store.New()
	s := "Hello, world"
	kv.Set(b("preset"), b(s))

	tests := []struct {
		start    int
		end      int
		expected string
	}{
		{0, 3, s[0:4]},
		{10, 100, s[10:]},
		{0, 0, s[:1]},
		{100, 10, ""},
		{100, 200, ""},
		{3, 11, s[3:]},
		{3, 12, s[3:]},
		{3, 14, s[3:]},
		{11, 11, s[11:]},
		{12, 12, ""},
		{3, -1, s[3:]},
		{-3, 1, ""},
		{-3, -1, s[9:]},
		{0, -1, s},
	}
	for _, tt := range tests {
		keys := map[string]string{"preset": tt.expected, "notset": ""}
		for k, v := range keys {
			input := []string{"GETRANGE", k, strconv.Itoa(tt.start), strconv.Itoa(tt.end)}
			got, _ := ExecuteCommand(kv, bA(input))
			if !reflect.DeepEqual(got, []byte(v)) {
				t.Errorf("ExecuteCommand(%q): got %q want %q", input, got, tt.expected)
			}
		}
	}
}

func Test__SETRANGE(t *testing.T) {
	kv := store.New()
	s := "Hello World"
	for i := 1; i < 6; i++ {
		kv.Set(b(fmt.Sprintf("preset%d", i)), b(s))
	}

	tests := []struct {
		key      string
		offset   int
		value    string
		expected string
	}{
		// offset < len(v); valLen < len(v)
		{"preset1", 6, "Redis", "Hello Redis"},
		// offset < len(v); valLen = len(v)
		{"preset2", 6, "Happy Seals", "Hello Happy Seals"},
		// offset < len(v); valLen < len(v)
		{"preset3", 0, "Happy", "Happy World"},
		// offset = len(v)
		{"preset4", 11, ", Bye", s + ", Bye"},
		// offset > len(v); valLen < len(v)
		{"preset5", len(s) + 1, "Bye", "Hello World\x00Bye"},
		// @ offset > len(v)
		{"notset1", 0, "Hello", "Hello"},
		// @ offset > len(v)
		{"notset2", 6, "Hello", "\x00\x00\x00\x00\x00\x00Hello"},
	}

	for _, tt := range tests {
		input := []string{"SETRANGE", tt.key, strconv.Itoa(tt.offset), tt.value}
		got, _ := ExecuteCommand(kv, bA(input))
		v, _ := kv.Get(b(tt.key))
		e := tt.expected
		eL := len(e)
		if (got != eL) || !bytes.Equal(v, b(e)) {
			t.Errorf("ExecuteCommand(%q): got %q with len %d want %q with len %d", input, v, got, e, eL)
		}
	}
}
