package corpus

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

func LoadDefaultTextAsWordlist() []string {
	textFile, err := Data.Open("data/loremipsum.txt")
	if err != nil {
		panic("cannot open embedded default corpus")
	}
	defer textFile.Close()

	result, err := LoadTextAsWordlist(textFile)
	if err != nil {
		panic(fmt.Errorf("cannot load embedded default corpus: %w", err))
	}

	return result
}

func LoadTextAsWordlist(r io.Reader) ([]string, error) {
	contentBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	content := string(contentBytes)
	if !utf8.ValidString(content) {
		return nil, fmt.Errorf("text is not a valid UTF-8 string")
	}

	return extractWords(content), nil
}

func extractWords(s string) []string {
	words := strings.Fields(s)
	for i, word := range words {
		word = strings.ToLower(word)
		parts := strings.FieldsFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		word = strings.Join(parts, "")
		word = strings.ReplaceAll(word, "ß", "ss")

		words[i] = word
	}

	return words
}
