// пакет тестов
package test

import (
	"common"
	"testing"
)

func TestValidString(t *testing.T) {
	tests := []struct {
		Input    string
		Expected bool
	}{
		{"", false},           // пустая строка
		{"   ", false},        // строка из пробелов
		{"he   llo", true},    // строка с пробелами внутри
		{"hello", true},       // строка без пробелов
		{"   hello   ", true}, // строка с пробелами в начале и конце
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			result := common.ValidString(tt.Input)
			if result != tt.Expected {
				t.Errorf("for input [%s]', expected [%t] but got [%t]", tt.Input, tt.Expected, result)
			}
		})
	}
}
