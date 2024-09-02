
# Package `godoc_readme`

<!-- THIS FILE IS GENERATED. DO NOT EDIT! -->
> [github.com/dubbikins/godoc-readme](https://github.com/dubbikins/godoc-readme)

Godoc-readme is a CLI that generates a README.md file for your go project using comments you already write!

In fact, this README.md file was generated using godoc-readme! :open_mouth:

godoc-readme is built using the [godoc](https://go.dev/blog/godoc) from the standard library.

Usage:
    godoc-readme [flags]

Flags:
    -h, --help              help for godoc-readme
    -r, --recursive         Recursively search for go packages in the directory and generate a README.md for each package (default true)
    -t, --template string   The template file to use for generating the README.md file

Installing:
    go install github.com/dubbikins/godoc-readme/godoc-readme

> [!TIP]
> Use the `//go:generate godoc-readme` directive in your module root to generate a README.md file for your packages when the `go generate` command is run.

---

Markdown Text Styling

| Style | Syntax | Example | Output |
| ------| ------ | ------- | ------ |
| Bold | `** **` OR `__ __` | `**This is bold text**` | **This is bold text** |
|Italic| `* *` OR `_ _` | `_This text is italicized_` | _This text is italicized_ |
| Strikethrough| `~~ ~~` | `~~This was mistaken text~~` | This was mistaken text |
| Bold and nested italic | `** **` with `_ _` | `**This text is _extremely_ important**` | **This text is _extremely_ important** |
| All bold and italic | `*** ***` | `***All this text is important***` | ***All this text is important*** |
| Subscript | `<sub> </sub>` | `This is a <sub>subscript</sub> text` | This is a <sub>subscript</sub> text |
| Superscript | `<sup> </sup>` | `This is a <sup>superscript</sup> text` | This is a <sup>superscript</sup> text |
| Comments | `<!-- -->` | `<!-- This content will not appear in the rendered Markdown -->` | (Nothing Gets Renders ;)) |

> [!TIP]
> Adding a `//go:generate godoc-readme` directive will generate a README.md file for your package when the `go generate` command is run.

Supported Github Markdown Features:

- [x] Headings
- [x] Alerts
- [x] Badges
- [x] Lists
  - [x] Nested Lists
- [x] Task Lists 😉
- [x] Images
- [x] Links
- [x] Tables
- [x] Code Blocks
- [x] Footnotes[^1]
  - [x] Multiline Footnotes[^2]
- [ ] Color Model

## Alerts

You can add [Github Markdown Alerts](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax#alerts) to your readme by utilizing the notes syntax in your godoc comments.
Godoc-readme support _in-line_ alerts in your ***packages*** godoc comments OR single-line alerts that "target" a the package or a `Type`, `Func`, `Method`, `Var`, or `Const`. Targeted types, besides a package, cannot have inlined alerts because their doc strings are nested by default, use targets for these types instead.
If you want to change this behaviour, you can provide your own template for custom rendering logic.

The following alert types are supported:

- `NOTE`
- `WARNING`
- `IMPORTANT`
- `CAUTION`
- `TIP`

Syntax:

    // TYPE(target): text

    Where `type` is one of the supported alert types and `target` is the name of the *package* or an exported *Type, Func, Method, Var, or Const in the package* that you want to target with the note.
    A single-line "targeted" Note will appear after the target's doc string section in the README.md file while in-line notes will appear in-line of the doc string.
    Targeted notes must be on a single line and must begin with a space.

> [!WARNING]
> An in-line alert cannot have whitespace before it's declaration or it will be rendered as plain doc string text while a targeted alert must have one space before it's declaration.

> [!TIP]
> In-line alerts are great for enhancing your documentaion in large godoc comments that you want to control the placement of the alert
while single-line alerts are great for adding a note to a specific type, func, method, var, or const in your package. Since a single-line alert doesn't have to be collocated with the target, you can add targeted alerts from anywhere in your package.

![Static Badge](https://img.shields.io/badge/build-passing-brightgreen)

<!-- Examples for footnotes-->
[^1]: A Footnote Example.
[^2]: To add line breaks within a footnote, prefix new lines with 2 spaces.
  This is a second line.

# Types

## [type PackageReadme](./readme.go#L197-L201)

>```go
>type PackageReadme struct {
>    Options ReadmeOptions
>    Pkg     *packages.Package
>    Doc     *doc.Package
>}
>```
>PackageReadme is a struct that holds the package, ast and docs of the package
>It's used to pass data to the readme template

--- 

## [type Readme](./readme.go#L102-L106)

>```go
>type Readme struct {
>    Pkgs    map[string]*packages.Package
>    pkgs    []*packages.Package
>    options *ReadmeOptions
>}
>```
>Readme is a struct that holds the packages, ast and docs of the package
>And is used to pass data to the readme template
>
>```mermaid
>classDiagram
>    class Readme
>    Readme : +map[string]*packages.Package Pkgs
>    Readme : -[]*packages.Package pkgs
>    Readme : -ReadmeOptions options
>    Readme : +Generate() error
>    Readme --> ReadmeOptions
>    Readme --> PackageReadme
>    class ReadmeOptions
>    ReadmeOptions : -string Dir
>    ReadmeOptions : -string DirPattern
>    ReadmeOptions : -string TemplateFile
>    class PackageReadme
>    PackageReadme : +ReadmeOptions Options
>    PackageReadme : +packages.Package Pkg
>    PackageReadme : +doc.Package Doc
>```

>[!NOTE]
>Because of the simpicity of godoc-readme's templating engine, you can add powerful customizations to your documentation like the class diagram that was created using a code block and the [mermaid.js](https://mermaid.js.org/) library that is supported out of the box with Github markdown. (not all features are supported though.)

---

### Methods

### [method Generate](./readme.go#L223-L294)

>```go
>func (readme *Readme) Generate() (err error)
>```
>Generate creates the README.md file for the packages that are registered with a `Readme`
>
>The README is generated in the directory of the package using the template file provided or the default template in none is provided.
>The following template functions available in the template engine are defined in the [`template_functions` package](./template_functions):
>| Function | Description | Example | Output |
>| --- | --- | --- | --- |
>| `example` | Renders a markdown representation of a `[doc.Example]` instance | `{{ example . }}` where `.` is a [doc.Example]| renders [an example like](/#Examples) |
>| `code` | Renders the start (or end) of a code block in markdown, optionally specifying the language format of the code block | `{{ code "go" }}fmt.Println("Hello World"){{ code }}` | `` ```go\nfmt.Println("Hello World")\n```\n`` |
>| `fmt` | Renders a formatted string representation of an [ast.Node] | `{{ fmt . }}` | `N/A` |
>| `link` | Renders a markdown link to the location of the [ast.Node] in a package | `{{ link "title" . }}` | `[title](...)` where ... is the relative link to the file ,including line numbers |
>| `alert` | Renders a markdown alert message based on the notes provided in the [doc.Package] | `{{ alert . "title" }}` | renders the alerts with the "title" target |
>| `section` | Renders an indented markdown section header | `{{ section "line 1 text\nline 2 text" 1}}` | `>line 1 text\n>line 2 text` |
>| `doc` | Renders a ***package's*** doc string, including in-line alerts, package ref's, etc | `{{ doc . }}` | `N/A` |
>| `relative_path` | Replaces the pwd the `.` | `{{ relative_path "/abs/path" }}` where `/abs` is the pwd | returns `./path` |
>
>Additionally, the following functions are available in the template engine:
>
>- `base`: [filepath.Base] Returns the base name of a file path
<details>
<summary>ExampleReadme_Generate</summary>

```go
func ExampleReadme_Generate{
    readme, err := NewReadme(func(ro *ReadmeOptions) {
        ro.Dir = "./examples/mermaid"
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
 // Output:
 // 
 // README.md file generated successfully :tada:
 // 
```

</details>

--- 

## [type ReadmeOptions](./readme.go#L110-L117)

>```go
>type ReadmeOptions struct {
>    Dir               string `env:"GODOC_README_MODULE_DIR"`
>    DirPattern        string `env:"GODOC_README_MODULE_DIR_PATTERN" default:"./..."`
>    TemplateFile      string `env:"GODOC_README_TEMPLATE_FILE"`
>    Format            func([]byte) []byte
>    package_load_mode packages.LoadMode
>}
>```
>ReadmeOptions is a struct that holds the options for the Readme struct
>You can set the options via the options functions or by setting the environment variables defined in the `env` struct tag for the Option field

--- 

## [type RenderFlag](./flags.go#L33-L33)

>```go
>type RenderFlag uint32
>```
> RenderFlags can be used to turn on and off rendering of different sections in the README.md file.
>
>The bitmask values are as follows:
>
>| 1 | 2 | 3 | 4 | 5 | 6 | 7 | ... | 32 |
>|---|---|---|---|---|---|---|-----|----|
>| Types | Funcs | TypeMethods | Vars | Consts | Examples | Alerts | TBD | RenderAll (default) |
>
>For example, to render only the types and functions in the README.md file, you would set the 1st and 2nd bits, i.e `0000 0011` or `RenderTypes | RenderFuncs`

---

### Methods

### [method IsSet](./flags.go#L36-L38)

>```go
>func (f RenderFlag) IsSet(flag RenderFlag) bool
>```
>IsSet returns true if the flag is set in the RenderFlags

--- 
---
# Functions

## [func Execute](./readme.go#L68-L76)

>```go
>func Execute(args ...string)
>```
>Execute runs the root command using the os.Args by default
>Optionally, you can pass in a list of arguments to run the command with

---
## [func FormatMarkdown](./readme.go#L185-L193)

>```go
>func FormatMarkdown(md []byte) []byte
>```
>FormatMarkdown applies the following formatting to the markdown:
>1. Replace all hard-tabs(`\t`) with 4 single space characters (`    `)
>2. Remove leading whitespace from blank lines
>3. Replace multiple `\n`(3+) with a single `\n`

---
## [func init](./readme.go#L34-L38)

>```go
>func init()
>```

---

## Vars
```go
// The readme templates are embedded in the binary so that it can be used as a default template
// This value can be overridden by providing a template file using the --template flag or the GODOC_README_TEMPLATE_FILE environment variable
//
//go:embed templates/*
var readme_templates embed.FS
```

```go
var recursive bool = true
```

```go
// The root command for the CLI
var rootCmd = &cobra.Command{
    Use:   "godoc-readme",
    Short: "Generate README.md file for your go project using comments you already write for godoc",
    Long:  `Generate README.md file for your go project using comments you already write for godoc`,
    Run: func(cmd *cobra.Command, args []string) {
        if readme, err := NewReadme(func(ro *ReadmeOptions) {
            if !recursive {
                ro.DirPattern = ro.Dir
            }
            if template_filename != "" {
                ro.TemplateFile = template_filename
            }
        }); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        } else {
            if err = readme.Generate(); err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }
        }

    },
}
```

```go
var template_file *os.File
```

```go
var template_filename string
```

# Examples

<details>
<summary>Example</summary>

```go
func Example{
    Execute()

}
 // Output:
 // 
 // README.md file generated successfully :tada:
 // 
```

</details>
<details>
<summary>Example_help_command</summary>

```go
func Example_help_command{
    Execute("-h")

}
 // Output:
 // 
 // Generate README.md file for your go project using comments you already write for godoc
 // 
 // Usage:
 //   godoc-readme [flags]
 // 
 // Flags:
 //   -h, --help   dhelp for godoc-readme
 //   -r, --recursive   Recursively search for go packages in the directory and generate a README.md for each package (default true)
 //   -t, --template string   The template file to use for generating the README.md file
 // 
```

</details>
<details>
<summary>Example_template_file</summary>

```go
func Example_template_file{
    Execute("-t", "README.tmpl")

}
 // Output:
 // 
 // README.md file generated successfully :tada:
 // 
```

</details>

## File Names

- [docs.go](./docs.go)
- [flags.go](./flags.go)
- [readme.go](./readme.go)

## Imports

- bytes
- embed
- fmt
- github.com/dubbikins/envy
- github.com/dubbikins/godoc-readme/template_functions
- github.com/spf13/cobra
- go/doc
- golang.org/x/tools/go/packages
- io
- log/slog
- os
- path
- path/filepath
- regexp
- strings
- text/template

