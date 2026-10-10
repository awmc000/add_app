// Package domain contains the pure, GUI-free logic for modelling a desktop
// launcher entry, serializing/parsing .desktop files, resolving icons, and
// installing to the user applications directory. It must never import Fyne.
package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// LauncherEntry is the set of fields collected by the form. See the
// Desktop Entry Specification (freedesktop.org) for the meaning of keys.
type LauncherEntry struct {
	Name        string
	Executable  string // absolute path or command name
	Arguments   string // optional, appended to Exec after the executable
	Comment     string // optional, shown as tooltip/subtitle
	Icon        string // optional; empty means "resolve automatically"
	InTerminal  bool
	Category    string // human label from the UI (see Categories)
}

// Categories maps the UI label to its freedesktop Categories value.
// An empty value means the Categories key is omitted.
var Categories = map[string]string{
	"General": "",
	"Development": "Development",
	"Game":        "Game",
	"Graphics":    "Graphics",
	"Network":     "Network",
	"Audio & Video": "AudioVideo",
	"Office":      "Office",
	"Science":     "Science",
	"System":      "System",
	"Utility":     "Utility",
}

// CategoryLabels returns the UI labels in display order.
func CategoryLabels() []string {
	return []string{
		"General", "Development", "Game", "Graphics", "Network",
		"Audio & Video", "Office", "Science", "System", "Utility",
	}
}

var unsafeFileChars = regexp.MustCompile(`[^a-z0-9]+`)

// FileName derives a safe, lowercase, dashed .desktop filename from Name.
// It returns an error only if nothing usable can be produced.
func FileName(name string) (string, error) {
	base := unsafeFileChars.ReplaceAllString(strings.ToLower(strings.TrimSpace(name)), "-")
	base = strings.Trim(base, "-")
	if base == "" {
		return "", errors.New("invalid name: cannot derive a filename")
	}
	return base + ".desktop", nil
}

// Validate checks that the required fields are present and usable.
func (e LauncherEntry) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("an application name is required")
	}
	if strings.TrimSpace(e.Executable) == "" {
		return errors.New("an executable is required")
	}
	if _, err := FileName(e.Name); err != nil {
		return err
	}
	return nil
}

// CategoryData returns the freedesktop value for the selected category label.
func (e LauncherEntry) CategoryData() string {
	return Categories[e.Category]
}

// DisplayExec returns the value used for the Exec= line: the executable
// followed by any arguments.
func (e LauncherEntry) DisplayExec() string {
	exe := strings.TrimSpace(e.Executable)
	args := strings.TrimSpace(e.Arguments)
	if args == "" {
		return exe
	}
	return fmt.Sprintf("%s %s", exe, args)
}
