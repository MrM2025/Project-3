package application 

import (
	"testing"
)

func TestCalculations(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		arg1, arg2      float64
		expected  float64
		expectErr bool
	}{
		{
			name:      "Addition positive numbers",
			operation: "+",
			arg1:         25.0,
			arg2:         3.5,
			expected:  28.5,
			expectErr: false,
		},
		{
			name:      "Addition negative numbers",
			operation: "+",
			arg1:         2.5,
			arg2:         -3.5,
			expected:  -1.0,
			expectErr: false,
		},

		{
			name:      "Subtraction positive numbers",
			operation: "-",
			arg1:         -2.0,
			arg2:         2.0,
			expected:  -4,
			expectErr: false,
		},
		{
			name:      "Subtraction negative numbers",
			operation: "-",
			arg1:         -5.0,
			arg2:         -2.5,
			expected:  -2.5,
			expectErr: false,
		},

		{
			name:      "Multiplication positive numbers",
			operation: "*",
			arg1:         2.0,
			arg2:         3.0,
			expected:  6.0,
			expectErr: false,
		},
		{
			name:      "Multiplication by zero",
			operation: "*",
			arg1:         1.0,
			arg2:         0.0,
			expected:  0.0,
			expectErr: false,
		},

		{
			name:      "Division positive numbers",
			operation: "/",
			arg1:         6.0,
			arg2:         2.0,
			expected:  3.0,
			expectErr: false,
		},
		{
			name:      "Division by zero",
			operation: "/",
			arg1:         5.0,
			arg2:         0.0,
			expected:  0.0,
			expectErr: true,
		},

		{
			name:      "Invalid operator",
			operation: "invalid",
			arg1:         2.0,
			arg2:         3.0,
			expected:  0.0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator(tt.operation, tt.arg1, tt.arg2, 10)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}