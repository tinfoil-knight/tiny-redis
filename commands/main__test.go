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
