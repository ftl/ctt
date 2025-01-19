package fyneui

import "fyne.io/fyne/v2"

type MainMenuController interface {
	LoadCorpus()
	Quit()
}

type mainMenu struct {
	controller MainMenuController

	root *fyne.MainMenu
	fileMenu
}

type fileMenu struct {
	fileLoadCorpus,
	fileQuit *fyne.MenuItem
}

func setupMainMenu(controller MainMenuController) *mainMenu {
	result := &mainMenu{
		controller: controller,
	}

	result.root = fyne.NewMainMenu(
		result.setupFileMenu(),
	)

	return result
}

// FILE

func (m *mainMenu) setupFileMenu() *fyne.Menu {
	m.fileLoadCorpus = fyne.NewMenuItem("Load Corpus...", m.controller.LoadCorpus)
	m.fileQuit = fyne.NewMenuItem("Quit", m.controller.Quit)
	m.fileQuit.IsQuit = true

	return fyne.NewMenu("File",
		m.fileLoadCorpus,
		fyne.NewMenuItemSeparator(),
		m.fileQuit,
	)
}

type nullMainMenuController struct{}

var _ MainMenuController = (*nullMainMenuController)(nil)

func (n *nullMainMenuController) LoadCorpus() {}

func (n *nullMainMenuController) Quit() {}
