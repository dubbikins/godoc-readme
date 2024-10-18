package template_functions

import (
	"bytes"
	"fmt"
	"go/doc"
	"go/format"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ExampleCode returns a function, given a package containing the example code, that returns a string representation of a doc.Example (Example Function in a package)
// You can call this function in a template by using `{{ example . }}` where `.` is a `*doc.Example` instance
func ExampleCode(pkg *packages.Package) func(*doc.Example) string {

	return func(ex *doc.Example) string {
		var buf = bytes.NewBuffer(nil)
		buf.WriteString("<details>\n")
		buf.WriteString(fmt.Sprintf("<summary>Example%s</summary>\n\n", ex.Name))
		buf.WriteString(fmt.Sprintf("```go\nfunc Example%s", ex.Name))
		format.Node(buf, pkg.Fset, ex.Code)
		output_lines := strings.Split(ex.Output, "\n")
		buf.WriteString("\n // Output:")
		for _, line := range output_lines {
			buf.WriteString(fmt.Sprintf("\n // %s", line))
		}
		buf.WriteString("\n```\n\n") // code blocks should be followed by a blank line
		buf.WriteString("</details>\n") 
		return buf.String()
	}
}
