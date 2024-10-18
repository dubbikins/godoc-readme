package template_functions

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var godoc_readme_badge_text = "![godoc-readme badge](https://img.shields.io/badge/generated%20by%20godoc--readme-00ADD8?style=plastic&logoSize=large&logo=Go&logoColor=00ADD8&labelColor=FFFFFF)\n"
func TestDocStringReplaceAllLeadingTabs(t *testing.T) {
	//Test that it only replaces leading tabs
	have := DocString("This is a test\n\twith a non-leading tab")
	want := "This is a test\n    with a non-leading tab"
	if have != want {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(have, want, true)
		t.Errorf("expected %q but got %q\nDiffs:\n%s", want, have, dmp.DiffPrettyText(diffs))
	}
}


func TestFormatTabs(t *testing.T) {
	have := PackageDocString("title\n## Alerts\n\nNOTE(target): this is a note\n\nNext line\n")
	want := godoc_readme_badge_text + "\n## Alerts\n\n> [!NOTE]\n> this is a note\n\nNext line\n"
	if have != want {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(have, want, true)
		t.Errorf("expected %q but got %q\nDiffs:\n%s", want, have, dmp.DiffPrettyText(diffs))
	}
}