package template_functions

import (
	"bytes"
	"regexp"
	"testing"
)

func TestXxx(t *testing.T) {
	var pattern = regexp.MustCompile(`(?m:^(NOTE|WARNING|IMPORTANT|CAUTION|TIP)\(([a-zA-Z][a-zA-Z0-9_]*)\):(.*)$)`)
	if !pattern.MatchString("NOTE(target): This is a note\ntes") {
		t.Errorf("does not match the expected pattern")
	}
	if !pattern.MatchString("WARNING(target): This is a note") {
		t.Errorf("does not match the expected pattern")
	}
	if !pattern.MatchString("IMPORTANT(target): This is a note") {
		t.Errorf("does not match the expected pattern")
	}
	if !pattern.MatchString("CAUTION(target): This is a note") {
		t.Errorf("does not match the expected pattern")
	}
	if !pattern.MatchString("TIP(target): This is a note") {
		t.Errorf("does not match the expected pattern")
	}
	if pattern.FindString("NOTE(target): This is a note\ntes") != "NOTE(target): This is a note" {
		t.Errorf("does not match the expected pattern")

	}

	// var replace = `>[!$1($2)]\n>$3`
	// var have = pattern.ReplaceAllString("NOTE(target): This is a note\ntes", replace)
	var submatches = pattern.FindSubmatch([]byte("NOTE(target): This is a note\ntes"))
	var expected_count = 8
	if len(submatches) != expected_count {
		t.Errorf("does not match the expected pattern.%s have %d, want %d", bytes.Join(submatches, []byte(";")), len(submatches), expected_count)
	}
}

func TestFormatInlineAlerts(t *testing.T) {

	have := DocString(`## Alerts

NOTE(target): this is a note

Next line
`)
	want := `## Alerts

> [!NOTE]
> this is a note

Next line
`
	if have != want {
		t.Errorf("expected %q but got %q", want, have)
	}
}
