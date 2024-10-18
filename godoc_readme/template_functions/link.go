package template_functions

import (
	"bytes"
	"fmt"
	"go/ast"
	"path"

	"golang.org/x/tools/go/packages"
)

// Link returns a markdown link to the  location of the ast.Node in a package
// Can be called in a template by using the `fmt` function `{{ link . }}` where `.` is a type that implements `*ast.Node`
func Link(pkg *packages.Package) func(string, ast.Node) string {
	return func(title string, node ast.Node) string {
		var buf = bytes.NewBuffer(nil)
		file := pkg.Fset.File(node.Pos())
		start_ln := file.Line(node.Pos())
		end_ln := file.Line(node.End())
		buf.WriteString(fmt.Sprintf("[%s](./%s#L%d-L%d)", title, path.Base(file.Name()), start_ln, end_ln))

		return buf.String()
	}
}
