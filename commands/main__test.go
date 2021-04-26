package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/tinfoil-knight/tiny-redis/store"
)

func Test__PING_ECHO(t *testing.T) {
	kv := store.New()
	tests := []struct {
		input    []string
		expected interface{}
	}{
		{[]string{"PING"}, "PONG"},
		{[]string{"PING", "hey"}, "hey"},
		{[]string{"ECHO", "hey"}, "hey"},
	}
	for _, tt := range tests {
		v, _ := ExecuteCommand(kv, tt.input)
		got := fmt.Sprintf("%s", v)
		if got != tt.expected {
			t.Errorf("ExecuteCommand(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}

func Test__GETRANGE(t *testing.T) {
	kv := store.New()
	s := "Hello, world"
	kv.Set("preset", s)

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
			got, _ := ExecuteCommand(kv, input)
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
		kv.Set(fmt.Sprintf("preset%d", i), s)
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
		got, _ := ExecuteCommand(kv, input)
		v, _ := kv.Get(tt.key)
		if got != len(tt.expected) || v != tt.expected {
			t.Errorf("ExecuteCommand(%q): got %q with len %d want %q with len %d", input, v, got, tt.expected, len(tt.expected))
		}
	}
}
