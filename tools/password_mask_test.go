package tools

import (
	"reflect"
	"testing"
)

func TestMaskPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no passwords",
			input:    []string{"ssh", "user@host", "--port", "22"},
			expected: []string{"ssh", "user@host", "--port", "22"},
		},
		{
			name:     "password flag with separate value",
			input:    []string{"mysql", "-password", "secret123", "-h", "localhost"},
			expected: []string{"mysql", "-password", "******", "-h", "localhost"},
		},
		{
			name:     "password flag with equals",
			input:    []string{"mysql", "--password=secret123", "-h", "localhost"},
			expected: []string{"mysql", "--password=******", "-h", "localhost"},
		},
		{
			name:     "multiple password flags",
			input:    []string{"app", "--api-key", "key123", "--token", "token456"},
			expected: []string{"app", "--api-key", "******", "--token", "******"},
		},
		{
			name:     "telegram token",
			input:    []string{"run4ever", "--notify-method", "telegram", "--telegram-token", "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"},
			expected: []string{"run4ever", "--notify-method", "telegram", "--telegram-token", "******"},
		},
		{
			name:     "password at end",
			input:    []string{"ssh", "-password", "secret"},
			expected: []string{"ssh", "-password", "******"},
		},
		{
			name:     "short flag p",
			input:    []string{"mysql", "-p", "secret123"},
			expected: []string{"mysql", "-p", "******"},
		},
		{
			name:     "auth token",
			input:    []string{"curl", "--auth-token", "bearer123", "https://api.example.com"},
			expected: []string{"curl", "--auth-token", "******", "https://api.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPassword(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MaskPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMaskPasswordEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty args",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "only password flag",
			input:    []string{"--password"},
			expected: []string{"--password"},
		},
		{
			name:     "password flag with equals and empty value",
			input:    []string{"--password="},
			expected: []string{"--password=******"},
		},
		{
			name:     "multiple equals in password",
			input:    []string{"--password=key=value=secret"},
			expected: []string{"--password=******"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPassword(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MaskPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}
