/*
Godoc-readme is a CLI that generates a README.md file for your go project using comments you already write!

In fact, this README.md file was generated using godoc-readme! :open_mouth:

godoc-readme is built using the [godoc](https://go.dev/blog/godoc) from the standard library.

Usage:

	godoc-readme [flags]

Flags:

	-h, --help              help for godoc-readme
	-r, --recursive         Recursively search for go packages in the directory and generate a README.md for each package (default true)
	-t, --template string   The template file to use for generating the README.md file

TIP(godoc_readme): Adding a `//go:generate godoc-readme` directive will generate a README.md file for your package when the `go generate` command is run.




Markdown Text Styling

| Style | Syntax | Example | Output |
| ------| ------ | ------- | ------ |
| Bold | `** **` OR `__ __` | `**This is bold text**` | **This is bold text** |
|Italic| `* *` OR `_ _` | `_This text is italicized_` | _This text is italicized_ |
| Strikethrough| `~~ ~~` | `~~This was mistaken text~~` | This was mistaken text |
| Bold and nested italic | `** **` with `_ _` | `**This text is _extremely_ important**` | **This text is _extremely_ important** |
| All bold and italic | `*** ***` | `***All this text is important***` | ***All this text is important*** |
| Subscript | `<sub> </sub>` | `This is a <sub>subscript</sub> text` | This is a <sub>subscript</sub> text |
| Superscript | `<sup> </sup> | `This is a <sup>superscript</sup> text` | This is a <sup>superscript</sup> text |

Supported Github Markdown Features:

- [x] Headings
- [x] Alerts
- [x] Lists
	- [x] Nested Lists
- [x] Task Lists 😉
- [x] Images
- [x] Links
- [x] Tables
- [x] Color Mode i.e. `#ffffff`
- [x] Code Blocks

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

WARNING(godoc_readme): An in-line alert cannot have whitespace before it's declaration or it will be rendered as plain doc string text while a targeted alert must have one space before it's declaration.

TIP(godoc_readme): In-line alerts are great for enhancing your documentaion in large godoc comments that you want to control the placement of the alert
while single-line alerts are great for adding a note to a specific type, func, method, var, or const in your package. Since a single-line alert doesn't have to be collocated with the target, you can add targeted alerts from anywhere in your package.

##
*/
package godoc_readme



/*
The standard templates provided for godoc-readme render a package's reading in the following format/order:

	1. Package Title (name)
	2. Package Docs
	3. Package Alerts
	4. Package Types (if any)

		1. Type w/Link to Source
		2. Type Signature
		3. Type Doc
		4. Type Alerts
		5. Type Methods
		6. Type Examples

	5. Package Functions (if any)
		1. Function Name w/Link to Source
		2. Function Signature
		3. Function Doc
		4. Function Alerts
		5. Function Examples

	6. Package Constants (if any)
		1. Constant Name w/Link to Source
		2. Constant Doc
	7. Package Variables (if any)
		1. Variables Name w/Link to Source
		2. Variables Doc
	8. Package Examples (if any)
		1. Example Name (dropdown)
		2. Example Code & Output
*/