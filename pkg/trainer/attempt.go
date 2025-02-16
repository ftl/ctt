package trainer

import "fmt"

type Attempt struct {
	CorrectPhrase string
	GivenPhrase   string
	Try           int
	Discarded     bool
}

func (a Attempt) Success() bool {
	return (a.CorrectPhrase == a.GivenPhrase)
}

func (a Attempt) String() string {
	result := a.GivenPhrase
	if result == "" {
		result = "(?)"
	}
	if a.Discarded {
		result = fmt.Sprintf("%s(discarded, %d)", a.CorrectPhrase, a.Try)
	} else if a.Try > 1 {
		result += fmt.Sprintf("(%d)", a.Try)
	}

	return result
}
