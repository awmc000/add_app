package app

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// openFilePicker opens a dedicated Fyne window file browser and calls onPick
// with the chosen absolute path (or onPick("") if the user cancels). filter,
// when non-empty, restricts selectable files to the listed extensions.
func openFilePicker(a fyne.App, title string, filter []string, onPick func(string)) {
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(920, 640))

	done := func(path string) {
		w.Close()
		onPick(path)
	}
	p := newPicker(w, filter, done, func() { w.Close() })
	w.SetContent(p.build())
	w.Show()
}

type entryInfo struct {
	name  string
	path  string
	isDir bool
}

// picker implements an in-app file browser.
type picker struct {
	w        fyne.Window
	cwd      string
	filter   []string
	onPick   func(string)
	onCancel func()

	path   *widget.Entry
	search *widget.Entry
	list   *widget.List
	open   *widget.Button

	entries  []entryInfo
	selected string
	content  fyne.CanvasObject
}

func newPicker(w fyne.Window, filter []string, onPick func(string), onCancel func()) *picker {
	p := &picker{
		w:        w,
		filter:   filter,
		onPick:   onPick,
		onCancel: onCancel,
		cwd:      homeDir(),
	}

	p.path = widget.NewEntry()
	p.path.SetText(p.cwd)
	p.path.OnSubmitted = func(v string) { p.enterPath(v) }

	p.search = widget.NewEntry()
	p.search.SetPlaceHolder("Search…")
	p.search.OnChanged = func(string) { p.refresh() }

	up := widget.NewButton("Up", func() { p.cd(filepath.Dir(p.cwd)) })
	home := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() { p.cd(homeDir()) })

	p.list = widget.NewList(
		func() int { return len(p.entries) },
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.FileIcon())
			label := widget.NewLabel("")
			return container.NewHBox(icon, label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			icon := box.Objects[0].(*widget.Icon)
			label := box.Objects[1].(*widget.Label)
			e := p.entries[id]
			if e.isDir {
				icon.SetResource(theme.FolderIcon())
			} else {
				icon.SetResource(theme.FileIcon())
			}
			label.SetText(e.name)
		},
	)
	p.list.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(p.entries) {
			return
		}
		e := p.entries[id]
		if e.isDir {
			p.list.UnselectAll()
			p.cd(e.path)
			return
		}
		p.selected = e.path
		p.open.Enable()
	}

	p.open = widget.NewButton("Open", func() {
		if p.selected != "" {
			p.onPick(p.selected)
		}
	})
	p.open.Importance = widget.HighImportance
	p.open.Disable()
	cancel := widget.NewButton("Cancel", p.onCancel)

	nav := container.NewBorder(nil, nil, container.NewHBox(up, home), nil, p.path)
	searchRow := container.NewBorder(nil, nil, widget.NewLabel("Search:"), nil, p.search)
	top := container.NewVBox(nav, searchRow)
	bottom := container.NewHBox(cancel, p.open)

	p.content = container.NewBorder(top, bottom, nil, nil, p.list)

	return p
}

func (p *picker) build() fyne.CanvasObject {
	p.refresh()
	return p.content
}

func (p *picker) cd(dir string) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return
	}
	p.cwd = dir
	p.path.SetText(dir)
	p.selected = ""
	p.open.Disable()
	p.refresh()
}

func (p *picker) enterPath(v string) {
	info, err := os.Stat(v)
	if err != nil {
		dialog.ShowError(err, p.w)
		return
	}
	if info.IsDir() {
		p.cd(v)
		return
	}
	p.selected = v
	p.onPick(v)
}

func (p *picker) refresh() {
	needle := strings.ToLower(strings.TrimSpace(p.search.Text))
	read, err := os.ReadDir(p.cwd)
	if err != nil {
		dialog.ShowError(err, p.w)
		return
	}

	var entries []entryInfo
	for _, de := range read {
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		isDir := de.IsDir()
		if !isDir && len(p.filter) > 0 && !hasExt(name, p.filter) {
			continue
		}
		if needle != "" && !strings.Contains(strings.ToLower(name), needle) {
			continue
		}
		entries = append(entries, entryInfo{name: name, path: filepath.Join(p.cwd, name), isDir: isDir})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].isDir != entries[j].isDir {
			return entries[i].isDir
		}
		return strings.ToLower(entries[i].name) < strings.ToLower(entries[j].name)
	})

	p.entries = entries
	p.list.UnselectAll()
	p.list.Refresh()
}

// homeDir returns the user's home directory, or the current working directory
// on failure (so the picker always has somewhere to start).
func homeDir() string {
	h, err := os.UserHomeDir()
	if err != nil || h == "" {
		if cwd, cerr := os.Getwd(); cerr == nil {
			return cwd
		}
		return "/"
	}
	return h
}

// hasExt reports whether name's extension (case-insensitive) is in the filter
// list, whose entries begin with a dot (e.g. ".png").
func hasExt(name string, filter []string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	for _, f := range filter {
		if strings.ToLower(f) == ext {
			return true
		}
	}
	return false
}
