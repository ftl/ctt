package fyneui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ftl/ctt/pkg/trainer"
)

type MainWindowController interface {
	Trainer
	Player
}

type Player interface {
	SetSpeed(int)
	SetFarnsworth(int)
	SetPitch(int)
}

type Trainer interface {
	Eval(string)
	SetMinLength(int)
	SetMaxLength(int)
	SetWordsPerPhrase(int)
}

type mainWindow struct {
	window  fyne.Window
	trainer Trainer
	player  Player

	input      *widget.Entry
	output     *widget.RichText
	outputText string

	speed          binding.Int
	farnsworth     binding.Int
	pitch          binding.Int
	minLength      binding.Int
	maxLength      binding.Int
	wordsPerPhrase binding.Int
}

func setupMainWindow(window fyne.Window, controller MainWindowController) *mainWindow {
	result := &mainWindow{
		window:  window,
		trainer: controller,
		player:  controller,
	}
	result.window.SetMaster()

	result.input = widget.NewEntry()
	result.input.OnSubmitted = result.inputSubmitted
	result.output = widget.NewRichText()
	result.output.Wrapping = fyne.TextWrapWord
	result.output.Scroll = 2 // widget.ScrollVerticalOnly

	result.speed = binding.NewInt()
	result.speed.AddListener(binding.NewDataListener(result.speedChanged))
	result.farnsworth = binding.NewInt()
	result.farnsworth.AddListener(binding.NewDataListener(result.farnsworthChanged))
	result.pitch = binding.NewInt()
	result.pitch.AddListener(binding.NewDataListener(result.pitchChanged))
	result.minLength = binding.NewInt()
	result.minLength.AddListener(binding.NewDataListener(result.minLengthChanged))
	result.maxLength = binding.NewInt()
	result.maxLength.AddListener(binding.NewDataListener(result.maxLengthChanged))
	result.wordsPerPhrase = binding.NewInt()
	result.wordsPerPhrase.AddListener(binding.NewDataListener(result.wordsPerPhraseChanged))

	root := container.NewBorder(
		container.NewGridWithColumns(2,
			container.NewBorder(nil, nil, widget.NewLabel("Speed:"), widget.NewLabel("WpM"), widget.NewEntryWithData(binding.IntToString(result.speed))),
			container.NewBorder(nil, nil, widget.NewLabel("Farnsworth:"), widget.NewLabel("WpM"), widget.NewEntryWithData(binding.IntToString(result.farnsworth))),
			container.NewBorder(nil, nil, widget.NewLabel("Pitch:"), widget.NewLabel("Hz"), widget.NewEntryWithData(binding.IntToString(result.pitch))),
			container.NewBorder(nil, nil, widget.NewLabel("Words per Phrase:"), nil, widget.NewEntryWithData(binding.IntToString(result.wordsPerPhrase))),
			container.NewBorder(nil, nil, widget.NewLabel("min. Length:"), widget.NewLabel("Characters"), widget.NewEntryWithData(binding.IntToString(result.minLength))),
			container.NewBorder(nil, nil, widget.NewLabel("max. Length:"), widget.NewLabel("Characters"), widget.NewEntryWithData(binding.IntToString(result.maxLength))),
		), // top
		container.NewBorder(nil, nil, widget.NewLabel("Input:"), nil, result.input), // bottom
		nil, // left
		nil, // right
		container.NewBorder(widget.NewLabel("Output:"), nil, nil, nil, result.output), // center
	)
	result.window.SetContent(root)

	return result
}

func (w *mainWindow) UseDefaultWindowGeometry() {
	w.window.Resize(fyne.NewSize(570, 300))
	w.window.CenterOnScreen()
}

func (w *mainWindow) SetMainMenu(menu *fyne.MainMenu) {
	w.window.SetMainMenu(menu)
	w.window.Canvas().Content().Refresh()
}

func (w *mainWindow) Show() {
	w.window.Show()
}

func (w *mainWindow) Reset() {
	w.outputText = ""
	w.output.ParseMarkdown(w.outputText)
	w.input.Text = ""
	w.window.Canvas().Focus(w.input)
}

func (w *mainWindow) Add(attempt trainer.Attempt) {
	md := attempt.GivenPhrase
	if attempt.Try > 1 {
		md += fmt.Sprintf("(%d)", attempt.Try)
	}
	if !attempt.Success() {
		md = fmt.Sprintf("*%s*", md)
	}
	md += " "
	w.outputText += md
	w.output.ParseMarkdown(w.outputText)
}

func (w *mainWindow) inputSubmitted(value string) {
	w.input.Text = ""
	w.input.Refresh()
	w.trainer.Eval(value)
}

func (w *mainWindow) SetSpeed(speed int) {
	w.speed.Set(speed)
}

func (w *mainWindow) speedChanged() {
	speed, _ := w.speed.Get()
	w.player.SetSpeed(speed)
}

func (w *mainWindow) farnsworthChanged() {
	fwpm, _ := w.farnsworth.Get()
	w.player.SetFarnsworth(fwpm)
}

func (w *mainWindow) SetPitch(pitch int) {
	w.pitch.Set(pitch)
}

func (w *mainWindow) pitchChanged() {
	pitch, _ := w.pitch.Get()
	w.player.SetPitch(pitch)
}

func (w *mainWindow) SetMinLength(minLength int) {
	w.minLength.Set(minLength)
}

func (w *mainWindow) minLengthChanged() {
	minLength, _ := w.minLength.Get()
	w.trainer.SetMinLength(minLength)
}

func (w *mainWindow) SetMaxLength(maxLength int) {
	w.maxLength.Set(maxLength)
}

func (w *mainWindow) maxLengthChanged() {
	maxLength, _ := w.maxLength.Get()
	w.trainer.SetMaxLength(maxLength)
}

func (w *mainWindow) SetWordsPerPhrase(wordsPerPhrase int) {
	w.wordsPerPhrase.Set(wordsPerPhrase)
}

func (w *mainWindow) wordsPerPhraseChanged() {
	wordsPerPhrase, _ := w.wordsPerPhrase.Get()
	w.trainer.SetWordsPerPhrase(wordsPerPhrase)
}

func (w *mainWindow) SelectOpenFile(callback func(string, error), title string, dir string, extensions ...string) {
	dirURI, err := storage.ListerForURI(storage.NewFileURI(dir))
	if err != nil {
		callback("", err)
		return
	}
	log.Printf("OPEN FILE in %s with extensions %v", dir, extensions)

	dialogCallback := func(r fyne.URIReadCloser, err error) {
		defer func() {
			if r != nil {
				r.Close()
			}
		}()
		if err != nil {
			callback("", err)
			return
		}
		if r == nil {
			callback("", nil)
			return
		}
		filename := r.URI().Path()
		log.Printf("file selected to open: %s", filename)
		callback(filename, nil)
	}

	fileDialog := dialog.NewFileOpen(dialogCallback, w.window)
	fileDialog.SetView(dialog.ListView)
	fileDialog.Resize(fyne.NewSize(1000, 600))
	// fileDialog.SetTitleText(title) // TODO: activate with fyne 2.6
	fileDialog.SetConfirmText("Open")
	fileDialog.SetDismissText("Cancel")
	fileDialog.SetLocation(dirURI)
	if len(extensions) > 0 {
		filterExtensions := make([]string, len(extensions), len(extensions))
		for i, extension := range extensions {
			filterExtensions[i] = "." + extension
		}
		fileDialog.SetFilter(storage.NewExtensionFileFilter(filterExtensions))
	}
	fileDialog.Show()
}

func (w *mainWindow) ShowInfoDialog(title string, format string, args ...any) {
	dialog.ShowInformation(
		title,
		fmt.Sprintf(format, args...),
		w.window,
	)
}

func (w *mainWindow) ShowErrorDialog(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	log.Println(err)
	dialog.ShowError(
		err,
		w.window,
	)
}
