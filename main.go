package main

import (
	appui "addapp/internal/app"
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
)

func main() {
	registerBundledIcons()

	a := fyneapp.NewWithID("dev.alex.addapp")
	w := a.NewWindow("Add App")
	w.Resize(fyne.NewSize(640, 560))

	appui.BuildUI(a, w)

	w.ShowAndRun()
}
