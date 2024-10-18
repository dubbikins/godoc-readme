package template_functions

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	"golang.org/x/tools/go/packages"
)



func GenDeclaration(pkg *packages.Package) func(*ast.GenDecl) string {

	return func(decl *ast.GenDecl) string {
		var buf = bytes.NewBuffer(nil)
		buf.WriteString("```go\n")
		var doc_ = decl.Doc
		decl.Doc = nil
		format.Node(buf, pkg.Fset, decl)
		decl.Doc = doc_
		buf.WriteString("\n```\n")
		return buf.String()
	}
}
type Target struct {
	start token.Pos // position of first character belonging to the node
    end token.Pos 
}
func (t *Target) Pos() token.Pos {
	return t.start
}
func (t *Target) End() token.Pos {
	return t.end
}

func FuncDeclaration(pkg *packages.Package) func(*ast.FuncDecl) string {

	return func(decl *ast.FuncDecl) string {
		var buf = bytes.NewBuffer(nil)
	
		buf.WriteString("```go\n")
		var doc = decl.Doc
		var body = decl.Body
		decl.Doc = nil
		decl.Body = nil
		format.Node(buf, pkg.Fset, decl)
		decl.Doc = doc
		decl.Body = body
		buf.WriteString("\n```\n")
		return buf.String()
	}
}

func SpecDeclaration(pkg *packages.Package) func([]ast.Spec) string {

	return func(decl []ast.Spec) string {
		var buf = bytes.NewBuffer(nil)
		buf.WriteString("```go\n")
		format.Node(buf, pkg.Fset, decl)
		buf.WriteString("\n```\n")
		return buf.String()
	}
}

func Declaration(pkg *packages.Package) func(ast.Node) string {

	return func(decl ast.Node) string {
		var buf = bytes.NewBuffer(nil)
		buf.WriteString("```go\n")
		format.Node(buf, pkg.Fset, decl)
		buf.WriteString("\n```\n")
		return buf.String()
	}
}