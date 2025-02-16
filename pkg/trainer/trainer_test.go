package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartWithEmptyPhrase(t *testing.T) {
	trainer, player, report, corpus := setupMockedTrainer()
	player.On("Play", "next").Once().Return()
	corpus.On("NextPhrase").Once().Return("next")

	trainer.Eval("")

	report.AssertNotCalled(t, "Add", mock.Anything)
	assert.Equal(t, "next", trainer.currentPhrase)
}

func TestCorrectValueIsReported(t *testing.T) {
	trainer, player, report, corpus := setupMockedTrainer()
	trainer.currentPhrase = "current"
	trainer.currentTry = 1
	correctAttempt := Attempt{CorrectPhrase: "current", GivenPhrase: "current", Try: 1}
	player.On("Play", "next").Once().Return()
	report.On("Add", correctAttempt).Once().Return()
	corpus.On("NextPhrase").Once().Return("next")

	trainer.Eval("current")

	report.AssertExpectations(t)
	assert.Equal(t, "next", trainer.currentPhrase)
}

func TestIncorrectValueIsReported(t *testing.T) {
	trainer, player, report, corpus := setupMockedTrainer()
	trainer.currentPhrase = "current"
	trainer.currentTry = 1
	correctAttempt := Attempt{CorrectPhrase: "current", GivenPhrase: "incorrect", Try: 1}
	player.On("Play", "current").Once().Return()
	report.On("Add", correctAttempt).Once().Return()

	trainer.Eval("incorrect")

	report.AssertExpectations(t)
	corpus.AssertNotCalled(t, "NextPhrase")
	assert.Equal(t, "current", trainer.currentPhrase)
}

func TestDiscardCurrentPhrase(t *testing.T) {
	trainer, player, report, corpus := setupMockedTrainer()
	trainer.currentPhrase = "current"
	trainer.currentTry = 1

	discardedAttempt := Attempt{CorrectPhrase: "current", Try: 1, Discarded: true}
	report.On("Add", discardedAttempt).Once().Return()
	corpus.On("NextPhrase").Once().Return("next")
	player.On("Play", "next").Once().Return()

	trainer.DiscardPhrase()

	report.AssertExpectations(t)
	assert.Equal(t, "next", trainer.currentPhrase)
}

// helpers

func setupMockedTrainer() (*Trainer, *mockPlayer, *mockReport, *mockCorpus) {
	player := &mockPlayer{}
	report := &mockReport{}
	corpus := &mockCorpus{}
	trainer := NewTrainer(player)
	trainer.SetReport(report)
	trainer.SetCorpus(corpus)
	return trainer, player, report, corpus
}

// mock types

var _ Player = (*mockPlayer)(nil)

type mockPlayer struct {
	mock.Mock
}

func (p *mockPlayer) Play(s string) {
	p.Called(s)
}

var _ Report = (*mockReport)(nil)

type mockReport struct {
	mock.Mock
}

func (m *mockReport) Add(r Attempt) {
	m.Called(r)
}

var _ Corpus = (*mockCorpus)(nil)

type mockCorpus struct {
	mock.Mock
}

func (m *mockCorpus) NextPhrase() string {
	args := m.Called()
	return args.String(0)
}
