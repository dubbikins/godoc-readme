package cmd

func Example_help_command() {
	Execute("-h")
	// Output:
	//
	// Generate README.md file for your go project using comments you already write
	//
	// Usage:
	//   godoc-readme [flags]
	//
	// Flags:
	//   -c, --confirm-updates   Use this flag to confirm overwriting existing README.md files. The default behaviour is to overwrite the file without confirmation. Confirmation also gives you the option to view the diff between the existing and generated file before overwriting it.
	//   -h, --help              help for godoc-readme
	//   -r, --recursive         If set, recursively search for go packages in the directory and generate a README.md for each package; Default will only create a Readme for the package found in the current directory
	//       --skip-examples     Skips generating the examples defined in test files
}

// func Example_template_file() {
// 	Execute("-t", "./templates/README.tmpl")
// 	// Output:
// 	//
// 	// Generated Documentation:
// 	// 	- ./templates/README.md
// }

// func Example() {

// 	err := Execute("-r")
// 	if err != nil {
// 		panic(err)
// 	}
// 	//	Output:
// 	//
// 	//	Generated Documentation:
// 	//	- ./README.md
// 	//	- ./godoc-readme/README.md
// 	//	- ./examples/mermaid/README.md
// 	//	- ./template_functions/README.md
// }

	 




