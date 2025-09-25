package instance

import (
	"strings"
	"testing"
)

func TestSafeLabel(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		labelType   LabelType
		expected    string
		expectError bool
	}{
		// Valid key tests
		{
			name:        "valid key - simple",
			input:       "valid-key",
			labelType:   LabelKey,
			expected:    "valid-key",
			expectError: false,
		},
		{
			name:        "valid key - with numbers",
			input:       "key-123",
			labelType:   LabelKey,
			expected:    "key-123",
			expectError: false,
		},
		{
			name:        "valid key - single letter",
			input:       "a",
			labelType:   LabelKey,
			expected:    "a",
			expectError: false,
		},

		// Key sanitization tests
		{
			name:        "key - starts with number",
			input:       "123test",
			labelType:   LabelKey,
			expected:    "n123test",
			expectError: false,
		},
		{
			name:        "key - replace slashes",
			input:       "test/key/name",
			labelType:   LabelKey,
			expected:    "test-key-name",
			expectError: false,
		},
		{
			name:        "key - replace underscores",
			input:       "test_key_name",
			labelType:   LabelKey,
			expected:    "test-key-name",
			expectError: false,
		},
		{
			name:        "key - replace colons",
			input:       "test:key:name",
			labelType:   LabelKey,
			expected:    "test-key-name",
			expectError: false,
		},
		{
			name:        "key - trim leading hyphen",
			input:       "-test-key",
			labelType:   LabelKey,
			expected:    "test-key",
			expectError: false,
		},
		{
			name:        "key - trim trailing hyphen",
			input:       "test-key-",
			labelType:   LabelKey,
			expected:    "test-key",
			expectError: false,
		},
		{
			name:        "key - trim both hyphens",
			input:       "-test-key-",
			labelType:   LabelKey,
			expected:    "test-key",
			expectError: false,
		},
		{
			name:        "key - too long",
			input:       strings.Repeat("a", 70),
			labelType:   LabelKey,
			expected:    strings.Repeat("a", 61),
			expectError: false,
		},
		{
			name:        "key - complex sanitization",
			input:       "123_test/key:name-",
			labelType:   LabelKey,
			expected:    "n123-test-key-name",
			expectError: false,
		},

		// Valid value tests
		{
			name:        "valid value - simple",
			input:       "valid-value",
			labelType:   LabelValue,
			expected:    "valid-value",
			expectError: false,
		},
		{
			name:        "valid value - starts with number",
			input:       "123test",
			labelType:   LabelValue,
			expected:    "123test",
			expectError: false,
		},
		{
			name:        "valid value - with underscores (converted to hyphens)",
			input:       "test_value",
			labelType:   LabelValue,
			expected:    "test-value",
			expectError: false,
		},

		// Value sanitization tests
		{
			name:        "value - replace slashes",
			input:       "test/value/name",
			labelType:   LabelValue,
			expected:    "test-value-name",
			expectError: false,
		},
		{
			name:        "value - replace underscores",
			input:       "test_value_name",
			labelType:   LabelValue,
			expected:    "test-value-name",
			expectError: false,
		},
		{
			name:        "value - replace colons",
			input:       "test:value:name",
			labelType:   LabelValue,
			expected:    "test-value-name",
			expectError: false,
		},
		{
			name:        "value - trim hyphens",
			input:       "-test-value-",
			labelType:   LabelValue,
			expected:    "test-value",
			expectError: false,
		},
		{
			name:        "value - too long",
			input:       strings.Repeat("a", 70),
			labelType:   LabelValue,
			expected:    strings.Repeat("a", 61),
			expectError: false,
		},

		// Error cases
		{
			name:        "empty string key",
			input:       "",
			labelType:   LabelKey,
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string value",
			input:       "",
			labelType:   LabelValue,
			expected:    "",
			expectError: true,
		},
		{
			name:        "key - only hyphens",
			input:       "---",
			labelType:   LabelKey,
			expected:    "",
			expectError: true,
		},
		{
			name:        "value - only hyphens",
			input:       "---",
			labelType:   LabelValue,
			expected:    "",
			expectError: true,
		},
		{
			name:        "key - invalid after sanitization",
			input:       "TEST",
			labelType:   LabelKey,
			expected:    "",
			expectError: true,
		},

		// Edge cases
		{
			name:        "key - single character after trim",
			input:       "-a-",
			labelType:   LabelKey,
			expected:    "a",
			expectError: false,
		},
		{
			name:        "value - single character after trim",
			input:       "-1-",
			labelType:   LabelValue,
			expected:    "1",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeLabel(tt.input, tt.labelType)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none. Result: %q", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

func TestLabelsValidate(t *testing.T) {
	tests := []struct {
		name        string
		labels      Labels
		expectError bool
	}{
		{
			name: "valid labels",
			labels: Labels{
				"valid-key":   "valid-value",
				"another-key": "123value",
			},
			expectError: false,
		},
		{
			name: "invalid key",
			labels: Labels{
				"123invalid": "valid-value",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.labels.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
