package domain

import (
	"bufio"
	"os"
	"strings"
)

// QuoteForExec wraps a token in double quotes if it contains whitespace,
// per the Desktop Entry Specification's Exec= quoting rules.
func QuoteForExec(token string) string {
	if strings.ContainsAny(token, " \t\n") {
		return `"` + strings.ReplaceAll(token, `"`, `\"`) + `"`
	}
	return token
}

// ExecLine builds a full Exec= value: a properly quoted executable plus args.
func ExecLine(e LauncherEntry) string {
	exe := QuoteForExec(strings.TrimSpace(e.Executable))
	if args := strings.TrimSpace(e.Arguments); args != "" {
		return exe + " " + args
	}
	return exe
}

// Serialize renders a LauncherEntry as .desktop file content, using the
// resolved icon path (already selected by the caller via ResolveIcon).
func Serialize(e LauncherEntry, iconPath string) string {
	var b strings.Builder
	b.WriteString("[Desktop Entry]\n")
	b.WriteString("Version=1.0\n")
	b.WriteString("Type=Application\n")
	b.WriteString("Name=" + e.Name + "\n")
	if comment := strings.TrimSpace(e.Comment); comment != "" {
		b.WriteString("Comment=" + comment + "\n")
	}
	b.WriteString("Exec=" + ExecLine(e) + "\n")
	if strings.TrimSpace(iconPath) != "" {
		b.WriteString("Icon=" + iconPath + "\n")
	}
	if e.InTerminal {
		b.WriteString("Terminal=true\n")
	} else {
		b.WriteString("Terminal=false\n")
	}
	if cat := e.CategoryData(); cat != "" {
		b.WriteString("Categories=" + cat + ";\n")
	}
	return b.String()
}

// ParsedEntry holds the keys we care about, read back from a .desktop file.
type ParsedEntry struct {
	LauncherEntry
}

// ParseDesktop reads .desktop content and populates a LauncherEntry.
// Unknown keys are ignored. Missing Exec/Name leave those fields empty so
// the UI can signal "no executable" gracefully.
func ParseDesktop(data string) (LauncherEntry, error) {
	var e LauncherEntry
	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		switch key {
		case "Name":
			e.Name = val
		case "Comment":
			e.Comment = val
		case "Terminal":
			e.InTerminal = strings.EqualFold(val, "true")
		case "Exec":
			SplitExec(&e, val)
		case "Icon":
			e.Icon = val
		case "Categories":
			for label, data := range Categories {
				if data != "" && strings.HasPrefix(val, data) {
					e.Category = label
					break
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		return e, err
	}
	return e, nil
}

// SplitExec extracts the executable path and remaining arguments from an
// Exec= line into e.Executable and e.Arguments.
func SplitExec(e *LauncherEntry, raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}

	// Pull off a quoted executable if present.
	if strings.HasPrefix(raw, `"`) {
		if end := strings.Index(raw[1:], `"`); end >= 0 {
			e.Executable = raw[1 : 1+end]
			rest := raw[end+2:]
			e.Arguments = strings.TrimSpace(unquoteArgs(rest))
			return
		}
	}

	// Otherwise the first whitespace-delimited token is the executable.
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return
	}
	e.Executable = strings.Trim(fields[0], `"'`)
	e.Arguments = strings.Join(fields[1:], " ")
}

// unquoteArgs strips trailing/leading quote noise from a remainder string.
func unquoteArgs(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// ParseDesktopFile reads a .desktop file from disk and parses it.
func ParseDesktopFile(path string) (LauncherEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LauncherEntry{}, err
	}
	return ParseDesktop(string(data))
}
