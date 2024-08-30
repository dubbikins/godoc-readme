/*
Godoc-readme is a CLI that generates a README.md file for your go project using comments you already write for godoc!

In fact, this README.md file was generated using godoc-readme! :open_mouth:

> [!Note]
> Adding a `//go:generate godoc-readme directive` to your go file will generate a README.md file for your package when the `go generate` command is run.

Usage:

	godoc-readme [flags]

Flags:

	-h, --help   help for auto-readme
	-r, --recursive   recursively search for go packages in the directory and generate a README.md for each package
*/
package godoc_readme
