package domain

import (
	"strings"
	"testing"
)

func TestFileName(t *testing.T) {
	cases := []struct {
		in   string
		want string
		err  bool
	}{
		{"My App", "my-app.desktop", false},
		{"Counter-Strike 2", "counter-strike-2.desktop", false},
		{"  GIMP   ", "gimp.desktop", false},
		{"a/b:c*d?e", "a-b-c-d-e.desktop", false},
		{"", "", true},
		{"   !!!   ", "", true},
	}
	for _, c := range cases {
		got, err := FileName(c.in)
		if c.err {
			if err == nil {
				t.Errorf("FileName(%q) expected error, got %q", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("FileName(%q) unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("FileName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestValidate(t *testing.T) {
	if err := (LauncherEntry{}).Validate(); err == nil {
		t.Error("expected error for empty entry")
	}
	if err := (LauncherEntry{Name: "X", Executable: ""}).Validate(); err == nil {
		t.Error("expected error for missing executable")
	}
	if err := (LauncherEntry{Name: "X", Executable: "/bin/true"}).Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestCategories(t *testing.T) {
	if got := (LauncherEntry{Category: "Game"}).CategoryData(); got != "Game" {
		t.Errorf("CategoryData(Game) = %q", got)
	}
	if got := (LauncherEntry{Category: "General"}).CategoryData(); got != "" {
		t.Errorf("CategoryData(General) = %q, want empty", got)
	}
	if got := (LauncherEntry{Category: "Audio & Video"}).CategoryData(); got != "AudioVideo" {
		t.Errorf("CategoryData(Audio & Video) = %q", got)
	}
}

func TestCategoryLabelsComplete(t *testing.T) {
	labels := CategoryLabels()
	for _, l := range labels {
		if _, ok := Categories[l]; !ok {
			t.Errorf("label %q missing from Categories map", l)
		}
	}
	// every categories value should be a valid semicolon token
	for l, v := range Categories {
		if v != "" && strings.ContainsAny(v, "; ") {
			t.Errorf("category %q has invalid value %q", l, v)
		}
	}
}
