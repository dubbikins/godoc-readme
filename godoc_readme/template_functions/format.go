package template_functions

import (
	"bytes"
	"go/ast"
	"go/format"
	"regexp"
	"strings"

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

// CAUTION(DocString): Targets types doc strings are nested by default and an alert will not be rendered correctly if they remain nested. If you are using the `DocString` function in a custom template setup, make sure you render the target's types without nesting to display the alerts correctly.

// PackageDocString returns a copy of *doc* with godoc notes replaced with github markdown notes
// Usage: `{{ DocString .Doc }}` where `.Doc` is a string containing godoc notes for a PACKAGE
func PackageDocString(doc string) string {

	var inline_alerts_pattern = regexp.MustCompile(`(?m:^(NOTE|WARNING|IMPORTANT|CAUTION|TIP)\(([a-zA-Z][a-zA-Z0-9_]*)\):(.*)$)`)
	var inline_alerts_replace = "> [!$1]\n>$3"
	var first_new_line_index = strings.IndexRune(doc, '\n')
	if first_new_line_index == -1 {
		first_new_line_index = 0
	}
	return DocString(inline_alerts_pattern.ReplaceAllString("![godoc-readme badge](https://img.shields.io/badge/generated%20by%20godoc--readme-00ADD8?style=plastic&logoSize=large&logo=Go&logoColor=00ADD8&labelColor=FFFFFF)\n" + doc[first_new_line_index:], inline_alerts_replace))

}

func DocString(doc string) string {
	var hard_tab_pattern = regexp.MustCompile(`(?m:^(\t+)(.*)$)`)
	hard_tab_replace_with_n_spaces := 4
	
	// var hard_tab_replace = fmt.Sprintf("%s$2", strings.Repeat(" ", hard_tab_replace_with_n_spaces))
	return string(hard_tab_pattern.ReplaceAllFunc([]byte( doc), func(b []byte) []byte {
		var replace = []byte{}
		for bytes.HasPrefix(b, []byte("\t")) {
			replace = append(replace, bytes.Repeat([]byte(" "), hard_tab_replace_with_n_spaces)...)
			b = b[1:]
		}
		return append(replace, b...)
	}))
}
