package app

import (
	"fmt"

	"addapp/internal/domain"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// form bundles all widget handles so handlers can read/write values.
type form struct {
	name      *widget.Entry
	executable *widget.Entry
	args      *widget.Entry
	comment   *widget.Entry
	icon      *widget.Entry
	iconPath  string
	terminal  *widget.Check
	category  *widget.Select
}

// controller owns the widget state and handles every user action.
type controller struct {
	a fyne.App
	w fyne.Window
	f form
}

// BuildUI assembles the window content and wires all callbacks.
func BuildUI(a fyne.App, w fyne.Window) {
	c := &controller{a: a, w: w}

	name := widget.NewEntry()
	name.SetPlaceHolder("My Application")

	executable := widget.NewEntry()
	executable.SetPlaceHolder("/path/to/executable or command")
	browseExec := widget.NewButton("Browse…", c.chooseExecutable)
	exeRow := container.NewBorder(nil, nil, nil, browseExec, executable)

	args := widget.NewEntry()
	args.SetPlaceHolder("--flag value  (optional)")

	comment := widget.NewEntry()
	comment.SetPlaceHolder("What does it do?  (optional)")

	icon := widget.NewEntry()
	icon.SetPlaceHolder("auto")
	icon.Disable()
	browseIcon := widget.NewButton("Browse…", c.chooseIcon)
	autoIcon := widget.NewButton("Auto", c.autoIcon)
	iconRow := container.NewBorder(nil, nil, nil, container.NewHBox(autoIcon, browseIcon), icon)

	terminal := widget.NewCheck("Run in terminal", nil)

	category := widget.NewSelect(domain.CategoryLabels(), nil)
	category.SetSelected("General")

	c.f = form{
		name:      name,
		executable: executable,
		args:      args,
		comment:   comment,
		icon:      icon,
		terminal:  terminal,
		category:  category,
	}

	frm := widget.NewForm(
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Executable", exeRow),
		widget.NewFormItem("Arguments", args),
		widget.NewFormItem("Comment", comment),
		widget.NewFormItem("Icon", iconRow),
		widget.NewFormItem("Terminal", terminal),
		widget.NewFormItem("Category", category),
	)

	save := widget.NewButton("Save .desktop", c.save)
	save.Importance = widget.HighImportance
	load := widget.NewButton("Load .desktop…", c.loadDesktop)
	quit := widget.NewButton("Quit", func() { w.Close() })

	buttons := container.NewHBox(save, load, quit)
	body := container.NewBorder(nil, buttons, nil, nil, container.NewVBox(frm))
	w.SetContent(container.NewPadded(body))
}

func (c *controller) currentEntry() domain.LauncherEntry {
	return domain.LauncherEntry{
		Name:       c.f.name.Text,
		Executable: c.f.executable.Text,
		Arguments:  c.f.args.Text,
		Comment:    c.f.comment.Text,
		Icon:       c.f.iconPath,
		InTerminal: c.f.terminal.Checked,
		Category:   c.f.category.Selected,
	}
}

func (c *controller) chooseExecutable() {
	openFileDialog(c.w, "Select executable", nil, func(path string) {
		if path == "" {
			return
		}
		c.f.executable.SetText(path)
		if c.f.iconPath == "" {
			c.refreshAutoIcon()
		}
	})
}

func (c *controller) chooseIcon() {
	openFileDialog(c.w, "Select icon image", domain.KnownIconExtensions, func(path string) {
		if path == "" {
			return
		}
		c.f.iconPath = path
		c.f.icon.SetText(path)
	})
}

func (c *controller) autoIcon() {
	c.f.iconPath = ""
	c.refreshAutoIcon()
}

// refreshAutoIcon re-runs icon resolution (exec-dir, then bundled default)
// and shows the result without touching an explicitly chosen icon.
func (c *controller) refreshAutoIcon() {
	entry := c.currentEntry()
	res := domain.ResolveIcon(entry, "")
	c.f.iconPath = res.Path
	c.f.icon.SetText(res.Path)
}

func (c *controller) save() {
	entry := c.currentEntry()
	if err := entry.Validate(); err != nil {
		dialog.ShowError(err, c.w)
		return
	}

	dir, err := domain.EnsureApplicationsDir()
	if err != nil {
		dialog.ShowError(err, c.w)
		return
	}
	target, err := domain.TargetPath(dir, entry.Name)
	if err != nil {
		dialog.ShowError(err, c.w)
		return
	}

	if !domain.Exists(target) {
		c.doWrite(entry, target)
		return
	}

	dialog.ShowConfirm(
		"Launcher already exists",
		fmt.Sprintf("A launcher already exists at:\n\n%s\n\nOverwrite it?\n(Use “Load .desktop…” to edit the existing file instead.)", target),
		func(overwrite bool) {
			if overwrite {
				c.doWrite(entry, target)
			}
		},
		c.w,
	)
}

func (c *controller) doWrite(entry domain.LauncherEntry, target string) {
	res := domain.ResolveIcon(entry, entry.Icon)
	content := domain.Serialize(entry, res.Path)
	if err := domain.WriteDesktop(content, target, true); err != nil {
		dialog.ShowError(err, c.w)
		return
	}
	dialog.ShowInformation("Done", "Created launcher:\n\n"+target, c.w)
	c.f.iconPath = res.Path
	c.f.icon.SetText(res.Path)
}

func (c *controller) loadDesktop() {
	openFileDialog(c.w, "Open .desktop to edit", []string{".desktop"}, func(path string) {
		if path == "" {
			return
		}
		entry, err := domain.ParseDesktopFile(path)
		if err != nil {
			dialog.ShowError(err, c.w)
			return
		}
		c.apply(entry)
	})
}

func (c *controller) apply(e domain.LauncherEntry) {
	c.f.name.SetText(e.Name)
	c.f.executable.SetText(e.Executable)
	c.f.args.SetText(e.Arguments)
	c.f.comment.SetText(e.Comment)
	c.f.terminal.SetChecked(e.InTerminal)
	if e.Category != "" {
		c.f.category.SetSelected(e.Category)
	}
	if e.Icon != "" {
		c.f.iconPath = e.Icon
		c.f.icon.SetText(e.Icon)
	} else {
		c.autoIcon()
	}
}
