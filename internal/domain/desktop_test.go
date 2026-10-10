package domain

import (
	"strings"
	"testing"
)

func TestExecLine(t *testing.T) {
	if got := ExecLine(LauncherEntry{Executable: "/bin/ls", Arguments: "-la"}); got != "/bin/ls -la" {
		t.Errorf("ExecLine = %q", got)
	}
	if got := ExecLine(LauncherEntry{Executable: "/home/alex/My App/bin"}); got != `"/home/alex/My App/bin"` {
		t.Errorf("ExecLine with spaces = %q", got)
	}
}

func TestSerializeBasics(t *testing.T) {
	e := LauncherEntry{
		Name:       "My App",
		Executable: "/bin/echo",
		Arguments:  "hi",
		Comment:    "says hi",
		InTerminal: true,
		Category:   "Utility",
	}
	got := Serialize(e, "/x/y/icon.png")
	for _, want := range []string{
		"[Desktop Entry]",
		"Type=Application",
		"Name=My App",
		"Comment=says hi",
		"Exec=/bin/echo hi",
		"Icon=/x/y/icon.png",
		"Terminal=true",
		"Categories=Utility;",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("serialize missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Icon=") && !strings.Contains(got, "Icon=/x/y/icon.png") {
		t.Error("unexpected Icon line")
	}
}

func TestSerializeOmitsOptionals(t *testing.T) {
	e := LauncherEntry{Name: "Min", Executable: "/bin/true", Category: "General"}
	got := Serialize(e, "")
	if strings.Contains(got, "Comment=") {
		t.Error("Comment should be omitted")
	}
	if strings.Contains(got, "Icon=") {
		t.Error("Icon should be omitted when empty")
	}
	if !strings.Contains(got, "Terminal=false") {
		t.Error("Terminal=false should be present")
	}
	if strings.Contains(got, "Categories=") {
		t.Error("Categories should be omitted for General")
	}
}

func TestParseDesktopRoundTrip(t *testing.T) {
	e := LauncherEntry{
		Name:       "My App",
		Executable: "/bin/echo",
		Arguments:  "hello world",
		Comment:    "a comment",
		Icon:       "/tmp/icon.png",
		InTerminal: true,
		Category:   "Development",
	}
	parsed, err := ParseDesktop(Serialize(e, e.Icon))
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Name != e.Name {
		t.Errorf("Name = %q", parsed.Name)
	}
	if parsed.Executable != e.Executable {
		t.Errorf("Executable = %q", parsed.Executable)
	}
	if parsed.Arguments != e.Arguments {
		t.Errorf("Arguments = %q", parsed.Arguments)
	}
	if parsed.Comment != e.Comment {
		t.Errorf("Comment = %q", parsed.Comment)
	}
	if parsed.Icon != e.Icon {
		t.Errorf("Icon = %q", parsed.Icon)
	}
	if !parsed.InTerminal {
		t.Error("InTerminal should be true")
	}
	if parsed.Category != e.Category {
		t.Errorf("Category = %q", parsed.Category)
	}
}

func TestParseDesktopQuotedExec(t *testing.T) {
	e, err := ParseDesktop("[Desktop Entry]\nType=Application\nName=Q\nExec=\"/path with space/app\" --flag\n")
	if err != nil {
		t.Fatal(err)
	}
	if e.Executable != "/path with space/app" {
		t.Errorf("Executable = %q", e.Executable)
	}
	if e.Arguments != "--flag" {
		t.Errorf("Arguments = %q", e.Arguments)
	}
}
