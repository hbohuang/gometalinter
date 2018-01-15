package api

import (
	"fmt"
	"strings"
)

// Severity of linter message
type Severity string

// Linter message severity levels.
const (
	Warning Severity = "warning"
	Error   Severity = "error"
)

func (s Severity) Less(other Severity) bool {
	return s != other && s == "warning"
}

func (s *Severity) UnmarshalText(text []byte) error {
	switch string(text) {
	case "warning":
		*s = Warning
	case "error":
		*s = Error
	default:
		return fmt.Errorf("invalid severity %q", string(text))
	}
	return nil
}

type Issue struct {
	Linter   string   `json:"linter"`
	Severity Severity `json:"severity"`
	Path     string   `json:"path"`
	Line     int      `json:"line"`
	Col      int      `json:"col"`
	Message  string   `json:"message"`
}

// NewIssue returns a new issue.
func NewIssue(linter string) *Issue {
	return &Issue{
		Line:     1,
		Severity: Warning,
		Linter:   linter,
	}
}

func (i *Issue) String() string {
	col := ""
	if i.Col != 0 {
		col = fmt.Sprintf("%d", i.Col)
	}
	return fmt.Sprintf("%s:%d:%s:%s: %s (%s)", strings.TrimSpace(i.Path), i.Line, col, i.Severity, strings.TrimSpace(i.Message), i.Linter)
}

// CompareIssue's Issues and return true if left should sort before right
func CompareIssue(l, r Issue, order []string) bool { // nolint: gocyclo
	for _, key := range order {
		switch {
		case key == "path" && l.Path != r.Path:
			return l.Path < r.Path
		case key == "line" && l.Line != r.Line:
			return l.Line < r.Line
		case key == "column" && l.Col != r.Col:
			return l.Col < r.Col
		case key == "severity" && l.Severity != r.Severity:
			return l.Severity < r.Severity
		case key == "message" && l.Message != r.Message:
			return l.Message < r.Message
		case key == "linter" && l.Linter != r.Linter:
			return l.Linter < r.Linter
		}
	}
	return true
}
