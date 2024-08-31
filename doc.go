/*
Godoc-readme is a CLI that generates a README.md file for your go project using comments you already write for godoc!

In fact, this README.md file was generated using godoc-readme! :open_mouth:

Usage:

	godoc-readme [flags]

Flags:

	-h, --help              help for godoc-readme
	-r, --recursive         Recursively search for go packages in the directory and generate a README.md for each package (default true)
	-t, --template string   The template file to use for generating the README.md file
*/
package godoc_readme

//NOTE(package-post-doc): Adding a `//go:generate godoc-readme directive` to your go file will generate a README.md file for your package when the `go generate` command is run.
