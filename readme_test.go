package godoc_readme

import (
	"fmt"
	"os"
)

func Example_help_command() {
	Execute("-h")
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

}

func Example_template_file() {
	Execute("-t", "README.tmpl")
	// Output:
	//
	// README.md file generated successfully :tada:
}

func Example() {
	Execute()
	// Output:
	//
	// README.md file generated successfully :tada:
}

func ExampleReadme_Generate() {
	readme, err := NewReadme(func(ro *ReadmeOptions) {
		ro.Dir = "example/nested"
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err = readme.Generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Output:
	//
	// README.md file generated successfully :tada:
}
