package template_functions

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
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
	var expected_count = 4
	if len(submatches) != expected_count {
		t.Errorf("does not match the expected pattern.%s have %d, want %d", bytes.Join(submatches, []byte(";")), len(submatches), expected_count)
	}
}

func TestFormatInlineAlerts(t *testing.T) {

	have := PackageDocString("title\n## Alerts\n\nNOTE(target): this is a note\n\nNext line\n")
	want := "![godoc-readme badge](https://img.shields.io/badge/generated%20by%20godoc--readme-00ADD8?style=plastic&logoSize=large&logo=Go&logoColor=00ADD8&labelColor=FFFFFF)\n\n## Alerts\n\n> [!NOTE]\n> this is a note\n\nNext line\n"
	if have != want {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(have, want, true)
		t.Errorf("expected %q but got %q\nDiffs:\n%s", want, have, dmp.DiffPrettyText(diffs))
	}
}
