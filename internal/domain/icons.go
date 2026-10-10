package domain

import (
	"os"
	"path/filepath"
	"strings"
)

// Icon resolution priority:
//   1. An explicitly user-chosen icon (strict preference).
//   2. A file beside the executable with the same base name and a known
//      image extension.
//   3. A bundled default for the selected category (provided by the caller
//      as bundledIcons[category] -> SVG bytes + extension).
//
// Bundled default icons are copied into the icons dir under a stable name
// derived from the app name, and the resulting absolute path is returned so
// GNOME renders it without icon-theme cache concerns.

// KnownIconExtensions are image extensions considered when probing the
// directory that contains the executable.
var KnownIconExtensions = []string{".png", ".svg", ".xpm", ".jpg", ".jpeg", ".ico"}

// IconResult describes how an icon was resolved.
type IconResult struct {
	Path  string // absolute path to use for Icon=, or "" if none
	Source string // one of: "explicit", "execdir", "default", "latest"
}

// BundledIcon is a default icon supplied by the embed layer.
type BundledIcon struct {
	Data []byte
	Ext  string // ".svg" or ".png"
}

// BundledIcons is a mapping from category label to embedded icon bytes.
var BundledIcons map[string]BundledIcon

// ResolveIcon computes the Icon= path for an entry. If userIcon is non-empty
// it wins. On any I/O error it degrades gracefully to whatever is available.
func ResolveIcon(e LauncherEntry, userIcon string) IconResult {
	if userIcon != "" {
		return IconResult{Path: userIcon, Source: "explicit"}
	}

	if fromDir := probeExecutableDir(e.Executable); fromDir != "" {
		return IconResult{Path: fromDir, Source: "execdir"}
	}

	if bin, ok := BundledIcons[e.Category]; ok && len(bin.Data) > 0 {
		iconDir, err := EnsureIconsDir()
		if err == nil {
			fn, _ := FileName(e.Name)
			base := strings.TrimSuffix(fn, ".desktop")
			dest := filepath.Join(iconDir, base+bin.Ext)
			_ = os.WriteFile(dest, bin.Data, 0o644)
			return IconResult{Path: dest, Source: "default"}
		}
	}

	return IconResult{Source: "latest"}
}

// probeExecutableDir looks for an image sharing the executable's base name
// in the same directory.
func probeExecutableDir(executable string) string {
	exe := strings.TrimSpace(executable)
	if exe == "" {
		return ""
	}
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	for _, k := range KnownIconExtensions {
		candidate := filepath.Join(dir, stem+k)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}
