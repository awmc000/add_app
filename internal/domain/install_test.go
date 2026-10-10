package domain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplicationsDirUsesXDG(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/tmp/datahome")
	dir, err := ApplicationsDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != filepath.Join("/tmp/datahome", "applications") {
		t.Errorf("dir = %q", dir)
	}
}

func TestEnsureApplicationsDirCreates(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_DATA_HOME", base)
	dir, err := EnsureApplicationsDir()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}

func TestWriteDesktopOverwriteSemantics(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "a.desktop")
	if err := WriteDesktop("one", target, false); err != nil {
		t.Fatal(err)
	}
	// second non-overwrite write must fail
	if err := WriteDesktop("two", target, false); err == nil {
		t.Error("expected ExistsError on non-overwrite write to existing file")
	}
	// overwrite true must succeed
	if err := WriteDesktop("two", target, true); err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "two" {
		t.Errorf("content = %q", data)
	}
	// executable bit must be set
	info, _ := os.Stat(target)
	if info.Mode().Perm()&0o111 == 0 {
		t.Error("file is not executable")
	}
}

func TestExists(t *testing.T) {
	p := filepath.Join(t.TempDir(), "f")
	if Exists(p) {
		t.Error("should not exist")
	}
	os.WriteFile(p, []byte("x"), 0o644)
	if !Exists(p) {
		t.Error("should exist")
	}
}
