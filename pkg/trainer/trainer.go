package trainer

import (
	"strings"
)

type Player interface {
	Play(string)
}

type Report interface {
	Add(Attempt)
}

type Corpus interface {
	NextPhrase() string
}

type Attempt struct {
	CorrectPhrase string
	GivenPhrase   string
	Try           int
}

func (r Attempt) Success() bool {
	return (r.CorrectPhrase == r.GivenPhrase)
}

type Trainer struct {
	player Player
	report Report
	corpus Corpus

	minLength int
	maxLength int

	currentPhrase string
	currentTry    int
}

func NewTrainer(player Player) *Trainer {
	return &Trainer{
		player: player,
		report: &nullReport{},
		corpus: &nullCorpus{},
	}
}

func (t *Trainer) SetReport(report Report) {
	t.report = report
}

func (t *Trainer) SetCorpus(corpus Corpus) {
	t.corpus = corpus
}

func (t *Trainer) Reset() {
	t.currentPhrase = ""
	t.currentTry = 0
}

func (t *Trainer) SetMinLength(minLength int) {
	t.minLength = minLength
}

func (t *Trainer) SetMaxLength(maxLength int) {
	t.maxLength = maxLength
}

func (t *Trainer) Eval(s string) {
	if t.starting() {
		t.Next()
		return
	}

	attempt := Attempt{
		CorrectPhrase: t.currentPhrase,
		GivenPhrase:   normalizePhrase(s),
		Try:           t.currentTry,
	}
	t.reportAttempt(attempt)

	if attempt.Success() {
		t.Next()
	} else {
		t.Repeat()
	}
}

func (t *Trainer) starting() bool {
	return t.currentPhrase == ""
}

func (t *Trainer) reportAttempt(attempt Attempt) {
	if attempt.CorrectPhrase == "" {
		return
	}

	t.report.Add(attempt)
}

func (t *Trainer) Next() {
	t.currentPhrase = t.pickNextPhrase()
	t.currentTry = 0
	t.Repeat()
}

func (t *Trainer) pickNextPhrase() string {
	result := ""
	nextPhraseOK := false
	tries := 0
	for !nextPhraseOK && (tries < 10) {
		result = normalizePhrase(t.corpus.NextPhrase())
		nextPhraseOK = ((t.minLength == 0) || (len(result) >= t.minLength)) &&
			((t.maxLength == 0) || (len(result) <= t.maxLength))
		tries += 1
	}
	return result
}

func normalizePhrase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (t *Trainer) Repeat() {
	t.currentTry += 1
	t.player.Play(t.currentPhrase)
}

type nullReport struct{}

var _ Report = (*nullReport)(nil)

func (n *nullReport) Add(Attempt) {}

type nullCorpus struct{}

var _ Corpus = (*nullCorpus)(nil)

func (n *nullCorpus) NextPhrase() string { return "" }
