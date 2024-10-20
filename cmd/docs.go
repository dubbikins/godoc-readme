/*
Godoc Readme CLI

You can use the godoc-readme CLI to generate a README.md file for your go project using comments you've already written! :open_mouth:

	Usage:
		godoc-readme [flags]

	Flags:
		-c, --confirm-updates   Use this flag to confirm overwriting existing README.md files.
			The default behaviour is to overwrite the file without confirmation.
			Confirmation also gives you the option to view the diff between the
			existing and generated file before overwriting it.
		-h, --help              help for godoc-readme
		-r, --recursive         If set, recursively search for go packages in the directory
			and generate a README.md for each package;
			Default behavior is to only create a README for the package in the current directory.
			--skip-all          Skips generating all sections besides the package documentation
			--skip-consts       Shows generating the consts section
			--skip-examples     Skips generating the examples section
			--skip-filenames    Skips generating the files section
			--skip-funcs        Skips generating the functions section
			--skip-imports      Skips generating the imports section
			--skip-vars         Skips generating the vars section
*/
package cmd
