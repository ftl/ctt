package fyneui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ftl/ctt/pkg/trainer"
)

type AppController interface {
	MainWindowController
	MainMenuController
}

type App struct {
	id         string
	version    string
	controller AppController

	app        fyne.App
	mainWindow *mainWindow
	mainMenu   *mainMenu
}

func NewApp(id string, version string, controller AppController) *App {
	result := &App{
		id:         id,
		version:    version,
		controller: controller,
	}

	result.app = app.NewWithID(id)
	result.app.Lifecycle().SetOnStarted(result.activate)

	return result
}

func (a *App) activate() {
	title := fmt.Sprintf("Clear Text Trainer %s", a.version)

	a.mainWindow = setupMainWindow(a.app.NewWindow(title), a.controller)
	a.mainWindow.UseDefaultWindowGeometry()
	a.mainWindow.Show()

	a.mainMenu = setupMainMenu(a.controller)
	a.mainWindow.SetMainMenu(a.mainMenu.root)

	a.mainWindow.SetSpeed(22)
	a.mainWindow.SetPitch(700)
	a.mainWindow.SetWordsPerPhrase(1)
	a.mainWindow.Reset()
}

func (a *App) Run() {
	a.app.Run()
}

func (a *App) Quit() {
	a.app.Quit()
}

func (a *App) Add(attempt trainer.Attempt) {
	a.mainWindow.Add(attempt)
}

func (a *App) SelectOpenFile(callback func(string, error), title string, dir string, extensions ...string) {
	a.mainWindow.SelectOpenFile(callback, title, dir, extensions...)
}

func (a *App) ShowInfoDialog(title string, format string, args ...any) {
	a.mainWindow.ShowInfoDialog(title, format, args...)
}

func (a *App) ShowErrorDialog(format string, args ...any) {
	a.mainWindow.ShowErrorDialog(format, args...)
}
