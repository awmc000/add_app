package main

import (
	"embed"
	"io/fs"

	"addapp/internal/domain"
)

//go:embed assets/icons/*.svg
var iconFS embed.FS

// registerBundledIcons populates domain.BundledIcons so icon auto-resolution
// can fall back to a bundled default per category.
func registerBundledIcons() {
	load := func(name string) domain.BundledIcon {
		data, err := fs.ReadFile(iconFS, "assets/icons/"+name)
		if err != nil {
			return domain.BundledIcon{}
		}
		return domain.BundledIcon{Data: data, Ext: ".svg"}
	}
	domain.BundledIcons = map[string]domain.BundledIcon{
		"General":      load("logo.svg"),
		"Development":  load("dev.svg"),
		"Game":         load("game.svg"),
		"Graphics":     load("graphics.svg"),
		"Network":      load("network.svg"),
		"Audio & Video": load("media.svg"),
		"Office":       load("office.svg"),
		"Science":      load("science.svg"),
		"System":       load("system.svg"),
		"Utility":      load("tool.svg"),
	}
}
