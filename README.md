# godoc_readme

Auto-Readme is a CLI that generates a README.md file for your go project using comments you already write for godoc!

Usage:

	auto-readme [flags]

Flags:

	-h, --help   help for auto-readme
	-r, --recursive   recursively search for go packages in the directory and generate a README.md for each package

> [!Note]
> Adding a `//go:generate auto-readme directive` to your go file will generate a README.md file for your package when the `go generate` command is run.

## Functions

### Execute

Execute runs the root command


#### Examples



