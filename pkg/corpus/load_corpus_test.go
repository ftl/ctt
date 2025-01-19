package corpus

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTestAsWordlist(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty file",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single word",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "multiple words",
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "words with punctuation",
			input:    "hello, world!",
			expected: []string{"hello", "world"},
		},
		{
			name:     "words with umlauts",
			input:    "Ärzte Ölscheich Übermut Straße",
			expected: []string{"ärzte", "ölscheich", "übermut", "strasse"},
		},
		{
			name:     "words with apostrophe",
			input:    "Futter'n wie bei Mutter'n", // kudos: http://www.deppenapostroph.info/wp-content/uploads/Futter-wie-bei-Muttern.jpg
			expected: []string{"futtern", "wie", "bei", "muttern"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(test.input)
			actual, err := LoadTextAsWordlist(reader)

			require.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestLoadTextAsWordlist_LoremIpsum(t *testing.T) {
	textFile, err := Data.Open("data/loremipsum.txt")
	require.NoError(t, err)
	defer textFile.Close()

	words, err := LoadTextAsWordlist(textFile)
	require.NoError(t, err)

	assert.Equal(t, 100, len(words))
}
