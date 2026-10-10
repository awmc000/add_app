package domain

import (
	"os"
	"path/filepath"
)

// Exists reports whether a path is present on disk.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// WriteDesktop writes serialized content to target, creating the parent
// directory if needed, and marks the file executable (required for GNOME to
// surface it in the app grid). overwrite must be true to replace an existing
// file; otherwise an existing file causes an error.
func WriteDesktop(content, target string, overwrite bool) error {
	if !overwrite && Exists(target) {
		return &ExistsError{Path: target}
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(target, []byte(content), 0o755); err != nil {
		return err
	}
	if err := os.Chmod(target, 0o755); err != nil {
		return err
	}
	return nil
}

// ExistsError is returned when a write target already exists.
type ExistsError struct {
	Path string
}

func (e *ExistsError) Error() string {
	return "already exists: " + e.Path
}
