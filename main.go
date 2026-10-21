package main

import (
	appui "addapp/internal/app"
	fyneapp "fyne.io/fyne/v2/app"
)

func main() {
	registerBundledIcons()

	a := fyneapp.NewWithID("dev.alex.addapp")
	w := a.NewWindow("Add App")

	content := appui.BuildUI(a, w)
	w.Resize(content.MinSize())

	w.ShowAndRun()
}
