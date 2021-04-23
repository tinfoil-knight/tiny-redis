package commands

import (
	"fmt"
	"testing"
)

func Test__PING_ECHO(t *testing.T) {
	tests := []struct {
		input    []string
		expected interface{}
	}{
		{[]string{"PING"}, "PONG"},
		{[]string{"PING", "hey"}, "hey"},
		{[]string{"ECHO", "hey"}, "hey"},
	}
	for _, tt := range tests {
		v, _ := ExecuteCommand(tt.input)
		got := fmt.Sprintf("%s", v)
		if got != tt.expected {
			t.Errorf("ExecuteCommand(%q): got %q want %q", tt.input, got, tt.expected)
		}
	}
}
