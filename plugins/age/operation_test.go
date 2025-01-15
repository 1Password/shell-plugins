package age

import "testing"

func TestOperationString(t *testing.T) {
	tests := []struct {
		name     string
		input    Operation
		expected string
	}{
		{
			name:     "Encrypt Operation",
			input:    Encrypt,
			expected: "encrypt",
		},
		{
			name:     "Decrypt Operation",
			input:    Decrypt,
			expected: "decrypt",
		},
		{
			name:     "Invalid Operation",
			input:    Operation(999),
			expected: "unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.String()
			if result != test.expected {
				t.Errorf("For input %v, expected %q but got %q", test.input, test.expected, result)
			}
		})
	}
}
