
# Package `godoc_readme`
<!-- THIS FILE IS GENERATED. DO NOT EDIT! -->
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

```shell
go install github.com/dubbikins/godoc-readme/godoc-readme
```

---

> [!TIP]
> Adding a `//go:generate godoc-readme` directive will generate a README.md file for your package when the `go generate` command is run.

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
Godoc-readme support *in-line* alerts in your _**packages**_ godoc comments OR single-line alerts that "target" a the package or a `Type`, `Func`, `Method`, `Var`, or `Const`. Targeted types, besides a package, cannot have inlined alerts because their doc strings are nested by default, use targets for these types instead.
If you want to change this behaviour, you can provide your own template for custom rendering logic.

The following alert types are supported:

  - NOTE
  - WARNING
  - IMPORTANT
  - CAUTION
  - TIP

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

<!-- Examples for footnotes-->`
[^1]: A Footnote Example.
[^2]: To add line breaks within a footnote, prefix new lines with 2 spaces.
  This is a second line.



## Types


### [type PackageReadme](./readme.go#L179-L179)

> ```go
>type PackageReadme struct {
>	Options ReadmeOptions
>	Pkg     *packages.Package
>	Doc     *doc.Package
>}
> ```
>PackageReadme is a struct that holds the package, ast and docs of the package
>It's used to pass data to the readme template
>



---





### [type Readme](./readme.go#L99-L99)

> ```go
>type Readme struct {
>	Pkgs map[string]*packages.Package
>	// contains filtered or unexported fields
>}
> ```
>Readme is a struct that holds the packages, ast and docs of the package
>And is used to pass data to the readme template
>
>```mermaid
>classDiagram
>	class Readme
>	Readme : +map[string]*packages.Package Pkgs
>    Readme : -[]*packages.Package pkgs
>	Readme : -ReadmeOptions options
>	Readme : +Generate() error
>	Readme --> ReadmeOptions
>	Readme --> PackageReadme
>	class ReadmeOptions
>	ReadmeOptions : -string Dir
>	ReadmeOptions : -string DirPattern
>	ReadmeOptions : -string TemplateFile
>	class PackageReadme
>	PackageReadme : +ReadmeOptions Options
>	PackageReadme : +packages.Package Pkg
>	PackageReadme : +doc.Package Doc
>```
>

>[!NOTE]
>Because of the simpicity of godoc-readme's templating engine, you can add powerful customizations to your documentation like the class diagram that was created using a code block and the [mermaid.js](https://mermaid.js.org/) library that is supported out of the box with Github markdown. (not all features are supported though.)



---

#### Methods

### [method Generate](./readme.go#L191-L191)
> ```go
> func (readme *Readme) Generate() (err error)
> ```
>Generate creates the README.md file for the packages that are registered with a `Readme`
>
>The README is generated in the directory of the package using the template file provided or the default template in none is provided.
>The template functions that are made available to the template arg defined in the [`template_functions` package](./template_functions)
>




### [type ReadmeOptions](./readme.go#L107-L107)

> ```go
>type ReadmeOptions struct {
>	Dir          string `env:"GODOC_README_MODULE_DIR"`
>	DirPattern   string `env:"GODOC_README_MODULE_DIR_PATTERN" default:"./..."`
>	TemplateFile string `env:"GODOC_README_TEMPLATE_FILE"`
>	// contains filtered or unexported fields
>}
> ```
>ReadmeOptions is a struct that holds the options for the Readme struct
>You can set the options via the options functions or by setting the environment variables defined in the `env` struct tag for the Option field
>



---




## Functions
### [func Execute](./readme.go#L65-L65)

>```go
>func Execute(args ...string)
>```
>
>Execute runs the root command using the os.Args by default
>Optionally, you can pass in a list of arguments to run the command with
>



---






## Examples

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




