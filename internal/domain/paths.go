package domain

import (
	"os"
	"path/filepath"
)

// ApplicationsDir returns the user's freedesktop applications directory,
// honoring XDG_DATA_HOME and defaulting to ~/.local/share/applications.
func ApplicationsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "applications"), nil
}

// EnsureApplicationsDir creates the applications directory if needed.
func EnsureApplicationsDir() (string, error) {
	dir, err := ApplicationsDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// TargetPath returns the full path a launcher would be written to for a given
// application name within the applications directory.
func TargetPath(dir, name string) (string, error) {
	fn, err := FileName(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fn), nil
}

// IconsDir returns where bundled/default icons are copied:
// $XDG_DATA_HOME/icons or ~/.local/share/icons.
func IconsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "icons"), nil
}

// EnsureIconsDir creates the icons directory if needed.
func EnsureIconsDir() (string, error) {
	dir, err := IconsDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
