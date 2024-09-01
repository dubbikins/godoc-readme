package template_functions

import (
	"bytes"
	"go/ast"
	"go/format"

	"golang.org/x/tools/go/packages"
)

// Format returns the string representation of an ast.Node in a package
// Can be called in a template by using the `fmt` function `{{ format . }}` where `.` is a type that implements `*ast.Node`
func FormatNode(pkg *packages.Package) func(ast.Node) string {

	return func(node ast.Node) string {
		var buf = bytes.NewBuffer(nil)
		if node != nil {
			format.Node(buf, pkg.Fset, node)
		}
		return buf.String()
	}
}
