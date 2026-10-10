package domain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveIconExplicitWins(t *testing.T) {
	res := ResolveIcon(LauncherEntry{Name: "A", Executable: "/bin/ls", Category: "General"}, "/user/picked.png")
	if res.Source != "explicit" || res.Path != "/user/picked.png" {
		t.Errorf("explicit resolution = %+v", res)
	}
}

func TestResolveIconExecDir(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "myapp")
	if err := os.WriteFile(exe, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "myapp.png"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	res := ResolveIcon(LauncherEntry{Name: "A", Executable: exe}, "")
	if res.Source != "execdir" {
		t.Fatalf("expected execdir, got %+v", res)
	}
	want := filepath.Join(dir, "myapp.png")
	if res.Path != want {
		t.Errorf("Path = %q, want %q", res.Path, want)
	}
}

func TestResolveIconBundledDefault(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	BundledIcons = map[string]BundledIcon{
		"Utility": {Data: []byte("<svg/>"), Ext: ".svg"},
	}
	res := ResolveIcon(LauncherEntry{Name: "My Tool", Executable: "/nonexistent/exe", Category: "Utility"}, "")
	if res.Source != "default" {
		t.Fatalf("expected default, got %+v", res)
	}
	if !filepath.IsAbs(res.Path) || filepath.Ext(res.Path) != ".svg" {
		t.Errorf("default icon path = %q", res.Path)
	}
	if !stringsContains(res.Path, "my-tool.svg") {
		t.Errorf("default icon path = %q, want to contain my-tool.svg", res.Path)
	}
}

func TestResolveIconNone(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	BundledIcons = nil
	res := ResolveIcon(LauncherEntry{Name: "X", Executable: "/nonexistent/exe"}, "")
	if res.Source != "latest" {
		t.Errorf("expected latest (none), got %+v", res)
	}
}

func stringsContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
