package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// EntryWithShortcut is a workaround until global shortcuts are also applied if
// an Entry widget is currently active.
// see https://github.com/fyne-io/fyne/issues/2627
// and https://github.com/fyne-io/fyne/issues/5355
type EntryWithShortcuts struct {
	widget.Entry

	shortcut fyne.ShortcutHandler
}

func NewEntryWithShortcuts() *EntryWithShortcuts {
	result := &EntryWithShortcuts{}
	result.ExtendBaseWidget(result)

	return result
}

func (e *EntryWithShortcuts) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	e.shortcut.AddShortcut(shortcut, handler)
}

func (e *EntryWithShortcuts) TypedShortcut(shortcut fyne.Shortcut) {
	if _, ok := shortcut.(*desktop.CustomShortcut); !ok {
		e.Entry.TypedShortcut(shortcut)
	}
	e.shortcut.TypedShortcut(shortcut)
}
