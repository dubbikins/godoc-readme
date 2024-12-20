package godoc_readme

import (
	"fmt"
	"os"
	"testing"
)


 


func TestFormatMarkdown(t *testing.T) {
	have := FormatMarkdown([]byte("\n\n"))
	want := "\n\n"
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}

	have = FormatMarkdown([]byte("\n  \t \n"))
	want = "\n\n"
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}

	have = FormatMarkdown([]byte("\n  \t \n\n"))
	want = "\n\n"
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}

	have = FormatMarkdown([]byte("\n\n\n"))
	want = "\n\n"
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}

	have = FormatMarkdown([]byte("\n\n\n\n\n\n\n\n\n\t    \n\n"))
	want = "\n\n"
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}

	have = FormatMarkdown([]byte("\t\t"))
	want = "        "
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}
	have = FormatMarkdown([]byte(`---


### [type Readme](./readme.go#L102-L102)`))
	want = `---

### [type Readme](./readme.go#L102-L102)`
	if string(have) != want {
		t.Errorf("have %q, want %q", have, want)
	}

}
func ExampleReadme_Generate() {
	readme, err := NewReadme(func(ro *ReadmeOptions) {
		ro.Dir = "../examples/mermaid"
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err = readme.Generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
}





