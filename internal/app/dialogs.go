// Package app contains the Fyne UI layer: window construction, native file
// dialogs, and the wiring that translates between widgets and the pure
// domain package.
package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

// openFileDialog launches a native file picker and calls onPick with the
// selected path, or onPick("") when the dialog is dismissed without a choice.
// filter, when non-empty, restricts visible files to the listed extensions.
func openFileDialog(parent fyne.Window, title string, filter []string, onPick func(string)) {
	fd := dialog.NewFileOpen(func(urc fyne.URIReadCloser, err error) {
		if err != nil {
			return
		}
		if urc == nil {
			onPick("")
			return
		}
		p := urc.URI().Path()
		_ = urc.Close()
		onPick(p)
	}, parent)
	if len(filter) > 0 {
		fd.SetFilter(storage.NewExtensionFileFilter(filter))
	}
	fd.Show()
}
