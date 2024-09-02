package template_functions

import (
	"os"
	"path"
	"testing"
)

func TestRelativeFilename(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	have := RelativeFilename(path.Join(pwd, "/main.go"))
	want := "./main.go"
	if have != want {
		t.Errorf("expected %q but got %q", want, have)
	}

	have = RelativeFilename(path.Join(pwd, "/path/to/dir"))
	want = "./path/to/dir"
	if have != want {
		t.Errorf("expected %q but got %q", want, have)
	}

	have = RelativeFilename(path.Join(pwd, "/nested/path/to/main.go"))
	want = "./nested/path/to/main.go"
	if have != want {
		t.Errorf("expected %q but got %q", want, have)
	}
}
