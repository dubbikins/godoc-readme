
# Package `template_functions`

<!-- THIS FILE IS GENERATED. DO NOT EDIT! -->
> [github.com/dubbikins/godoc-readme/template_functions](https://github.com/dubbikins/godoc-readme/template_functions)

Package template_functions provides a set of functions that are passed to the template engine to generate documentation.

The functions are used to format the documentation in a way that is easy to read and understand.

You can utilize these functions in your own custom templates to generate documentation for your packages with customize formatting/behvior if the standard templates provided by godoc-readme do not meet your needs.

## File Names

- ./template_functions/alert.go
- ./template_functions/code.go
- ./template_functions/decl.go
- ./template_functions/doc.go
- ./template_functions/example.go
- ./template_functions/filenames.go
- ./template_functions/format.go
- ./template_functions/link.go
- ./template_functions/section.go

## Imports

- bytes
- fmt
- go/ast
- go/doc
- go/format
- go/token
- golang.org/x/tools/go/packages
- os
- path
- regexp
- strings
- testing

# Types

## [type Target](./decl.go#L27-L27)

>```go
>type Target struct {
>    start token.Pos // position of first character belonging to the node
>    end   token.Pos
>}
>```

---

### Methods

### [method End](./decl.go#L34-L34)

>```go
>func (t *Target) End() token.Pos
>```

### [method Pos](./decl.go#L31-L31)

>```go
>func (t *Target) Pos() token.Pos
>```

--- 
---
# Functions

## [func Alert](./alert.go#L24-L24)

>```go
>func Alert(pkg *packages.Package, notes map[string][]*doc.Note) func(string) string
>```
>Alert returns a function that, given the name of a target, returns a string representing the alerts for that target
>Can be used in a template by calling `{{ Alert "target_name" }}` where `target_name` is the name of the package, a Type, Func, Var, or Const in the package.
>Alerts are rendered AFTER the doc comment for the target by default. Provide your own templates to modify this behavior.

>[!NOTE]
>Use this alert to provide additional information

>[!WARNING]
>Use this alert to warn users about potential issues

>[!IMPORTANT]
>Use this alert to highlight important information

>[!CAUTION]
>Use this alert to caution users about serious issues

>[!TIP]
>Use this alert to provide helpful tips to users

---
## [func CodeBlock](./code.go#L38-L38)

>```go
>func CodeBlock(pkg *packages.Package) func(lang ...string) string
>```
>CodeBlock returns the start (or end) of a code block in markdown
>If you provide a language, it will be used to specify the language format of the code block
>In a template, you should call this function like this: `{{ code_block "go" }}` or `{{ code_block }}` to start the code block
>and `{{ code_block }}` to end the code block after you've rendered the contens between the start and end of the code block
>
>Example:
>
>```tmpl
>{{ code_block "go" }}
>package main
>
>    func main() {
>        fmt.Println("Hello, World!")
>    }
>
>{{ code_block }}
>```
>This will render the following markdown:
>
>```go
>package main
>
>    func main() {
>        fmt.Println("Hello, World!")
>    }
>
>```

---
## [func DocString](./format.go#L38-L38)

>```go
>func DocString(doc string) string
>```

>[!CAUTION]
>Targets types doc strings are nested by default and an alert will not be rendered correctly if they remain nested. If you are using the `DocString` function in a custom template setup, make sure you render the target's types without nesting to display the alerts correctly.

---
## [func ExampleCode](./example.go#L15-L15)

>```go
>func ExampleCode(pkg *packages.Package) func(*doc.Example) string
>```
>ExampleCode returns a function, given a package containing the example code, that returns a string representation of a doc.Example (Example Function in a package)
>You can call this function in a template by using `{{ example . }}` where `.` is a `*doc.Example` instance

---
## [func FormatNode](./format.go#L15-L15)

>```go
>func FormatNode(pkg *packages.Package) func(ast.Node) string
>```
>Format returns the string representation of an ast.Node in a package
>Can be called in a template by using the `fmt` function `{{ format . }}` where `.` is a type that implements `*ast.Node`

---
## [func FuncDeclaration](./decl.go#L38-L38)

>```go
>func FuncDeclaration(pkg *packages.Package) func(*ast.FuncDecl) string
>```

---
## [func GenDeclaration](./decl.go#L14-L14)

>```go
>func GenDeclaration(pkg *packages.Package) func(*ast.GenDecl) string
>```

---
## [func Link](./link.go#L14-L14)

>```go
>func Link(pkg *packages.Package) func(string, ast.Node) string
>```
>Link returns a markdown link to the  location of the ast.Node in a package
>Can be called in a template by using the `fmt` function `{{ link . }}` where `.` is a type that implements `*ast.Node`

---
## [func PackageDocString](./format.go#L30-L30)

>```go
>func PackageDocString(doc string) string
>```
>PackageDocString returns a copy of *doc* with godoc notes replaced with github markdown notes
>Usage: `{{ DocString .Doc }}` where `.Doc` is a string containing godoc notes for a PACKAGE

---
## [func RelativeFilename](./filenames.go#L8-L8)

>```go
>func RelativeFilename(abs string) (relative string)
>```

---
## [func Section](./section.go#L10-L10)

>```go
>func Section(doc string, n int) string
>```
>Section returns a copy of doc with `n` number of `>`'s added to the beginning of each line.
>You can call this function in a template by using `{{ section .Doc 1 }}` where `.Doc` is *string* field.
>Example:
>
>    `Section("This is a section", 1)` returns "> This is a section"

---
## [func SpecDeclaration](./decl.go#L56-L56)

>```go
>func SpecDeclaration(pkg *packages.Package) func([]ast.Spec) string
>```

---
## [func TestFormatTabs](./format.go#L53-L53)

>```go
>func TestFormatTabs(t *testing.T)
>```

---

