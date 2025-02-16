package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttemptToString(t *testing.T) {
	tests := []struct {
		name     string
		input    Attempt
		expected string
	}{
		{
			name: "success on first try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "current",
				Try:           1,
			},
			expected: "current",
		},
		{
			name: "success on second try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "current",
				Try:           2,
			},
			expected: "current(2)",
		},
		{
			name: "no idea on first try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "",
				Try:           1,
			},
			expected: "(?)",
		},
		{
			name: "no idea on second try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "",
				Try:           2,
			},
			expected: "(?)(2)",
		},
		{
			name: "failure on first try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "cur",
				Try:           1,
			},
			expected: "cur",
		},
		{
			name: "failure on second try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "cur",
				Try:           2,
			},
			expected: "cur(2)",
		},
		{
			name: "discarded on second try",
			input: Attempt{
				CorrectPhrase: "current",
				GivenPhrase:   "",
				Try:           2,
				Discarded:     true,
			},
			expected: "current(discarded, 2)",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.input.String()
			assert.Equal(t, test.expected, actual)
		})
	}
}
