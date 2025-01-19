package app

import (
	"os"

	"github.com/ftl/ctt/pkg/corpus"
	"github.com/ftl/ctt/pkg/player"
	"github.com/ftl/ctt/pkg/trainer"
)

type Quitter interface {
	Quit()
}

type UI interface {
	SelectOpenFile(callback func(filename string, err error), title string, dir string, extensions ...string)
	ShowInfoDialog(title string, format string, args ...any)
	ShowErrorDialog(string, ...any)
}

type App struct {
	*trainer.Trainer
	*player.PAPlayer
	Quitter Quitter
	UI      UI

	wordlist []string
}

func NewApp() *App {
	result := &App{
		Quitter: &nullQuitter{},
		UI:      &nullUI{},
	}

	result.PAPlayer = player.NewPAPlayer(22, 770)
	result.Trainer = trainer.NewTrainer(result.PAPlayer)

	result.wordlist = corpus.LoadDefaultTextAsWordlist()
	result.SetCorpus(corpus.NewRandomCorpus(result.wordlist...))

	return result
}

func (a *App) LoadCorpus() {
	a.UI.SelectOpenFile(a.loadCorpus, "Load Corpus", "", "txt")
}

func (a *App) loadCorpus(filename string, err error) {
	if err != nil {
		a.UI.ShowErrorDialog("Cannot select a corpus file: %v", err)
		return
	}
	if filename == "" {
		return
	}

	text, err := os.Open(filename)
	if err != nil {
		a.UI.ShowErrorDialog("Cannot open the corpus file: %v", err)
		return
	}
	defer text.Close()

	wordlist, err := corpus.LoadTextAsWordlist(text)
	if err != nil {
		a.UI.ShowErrorDialog("Cannot load the corpus word list: %v", err)
		return
	}
	a.wordlist = wordlist
	corpus := corpus.NewRandomCorpus(a.wordlist...)
	a.Trainer.SetCorpus(corpus)
	a.Trainer.Reset()
}

func (a *App) Quit() {
	a.PAPlayer.SetSpeed(35)
	a.PAPlayer.Play("73")
	a.Quitter.Quit()
}

type nullQuitter struct{}

var _ Quitter = (*nullQuitter)(nil)

func (n *nullQuitter) Quit() {}

type nullUI struct{}

var _ UI = (*nullUI)(nil)

func (n *nullUI) SelectOpenFile(callback func(filename string, err error), title string, dir string, extensions ...string) {
}

func (n *nullUI) ShowInfoDialog(title string, format string, args ...any) {}

func (n *nullUI) ShowErrorDialog(string, ...any) {}
