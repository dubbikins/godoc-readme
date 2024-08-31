// more details about the package
package godoc_readme

//go:generate go run ./cmd/main.go
import (
	"bytes"
	_ "embed"
	"fmt"
	"go/doc"
	"go/format"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dubbikins/envy"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

// the readme_template is embedded in the binary so that it can be used as a default template
// This value can be overridden by providing a template file using the --template flag or the GODOC_README_TEMPLATE_FILE environment variable
//
//go:embed README.tmpl
var readme_template string
var recursive bool = true
var template_filename string
var template_file *os.File

func init() {
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Recursively search for go packages in the directory and generate a README.md for each package")
	rootCmd.PersistentFlags().StringVarP(&template_filename, "template", "t", "", "The template file to use for generating the README.md file")
	//rootCmd.Flags().BoolP("recursive", "r", true, "Recursively search for go packages in the directory and generate a README.md for each package")
}

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
		fmt.Println("README.md file generated successfully :tada:")
	},
}

// Execute runs the root command using the os.Args by default
// Optionally, you can pass in a list of arguments to run the command with
func Execute(args ...string) {
	if len(args) > 0 {
		rootCmd.SetArgs(args)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//Note(PackageReadme): test note

// Readme is a struct that holds the packages, ast and docs of the package
// And is used to pass data to the readme template
/*
```mermaid
classDiagram
	note "From Duck till Zebra"
    PackageReadme <|-- Duck
    note for Duck "can fly\ncan swim\ncan dive\ncan help in debugging"
    PackageReadme <|-- Fish
    PackageReadme <|-- Zebra
    PackageReadme : +int age
    PackageReadme : +String gender
    PackageReadme: +isMammal()
    PackageReadme: +mate()

```

*/
type Readme struct {
	pkgs        []*packages.Package
	RefinedPkgs map[string]*packages.Package
	Docs        []*doc.Package
	options     *ReadmeOptions
}

// ReadmeOptions is a struct that holds the options for the Readme struct
// You can set the options via the options functions or by setting the environment variables defined in the `env` struct tag for the Option field
type ReadmeOptions struct {
	Dir               string `env:"GODOC_README_MODULE_DIR"`
	DirPattern        string `env:"GODOC_README_MODULE_DIR_PATTERN" default:"./..."`
	TemplateFile      string `env:"GODOC_README_TEMPLATE_FILE"`
	package_load_mode packages.LoadMode
}

// NewReadme creates a new Readme struct that holds the packages, ast and docs of the package
// It loads the packages in the directory provided and parses the ast and docs of the package
// It returns an error if the packages failed to load
// The options are provided as functional options
// Usage:
// ```go
//
//	readme, err := NewReadme(func(ro *ReadmeOptions) {
//		ro.Dir = "./path/to/package"
//		ro.DirPattern = "./path/to/package/..."
//	})
//
// ```
func NewReadme(opts ...func(*ReadmeOptions)) (readme *Readme, err error) {
	readme = &Readme{
		options: &ReadmeOptions{
			package_load_mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypes | packages.NeedEmbedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedDeps,
		},
		RefinedPkgs: map[string]*packages.Package{},
	}
	if err = envy.Unmarshal(readme.options); err != nil {
		return
	}
	for _, opt := range opts {
		opt(readme.options)
	}

	cfg := &packages.Config{
		Mode:  readme.options.package_load_mode,
		Dir:   readme.options.Dir,
		Tests: true,
	}
	readme.pkgs, err = packages.Load(cfg, readme.options.DirPattern)
	if err != nil {
		return
	}
	if packages.PrintErrors(readme.pkgs) > 0 {
		// Package failed to parse
		os.Exit(1)
	}

	readme.Docs = []*doc.Package{}

	for _, pkg := range readme.pkgs {
		for f_index, file := range pkg.CompiledGoFiles {
			//Strip all of the files that are not go files
			if !strings.HasSuffix(file, ".go") {
				pkg.Syntax = append(pkg.Syntax[:f_index], pkg.Syntax[f_index+1:]...)
			}
		}
		if len(pkg.Syntax) == 0 {
			continue
		}
		//If the package is already in the refined packages map, then we check if the package is not a test package
		// The test package will include more details that we want
		if p, found := readme.RefinedPkgs[pkg.Name]; found {
			if !strings.Contains(p.ID, "test") {
				readme.RefinedPkgs[pkg.Name] = pkg
			}
		}
		readme.RefinedPkgs[pkg.Name] = pkg
	}
	return
}

// PackageReadme is a struct that holds the package, ast and docs of the package
// It's used to pass data to the readme template
type PackageReadme struct {
	Options ReadmeOptions
	P       *packages.Package
	Pkg     *doc.Package
}

// ExampleCode returns a function that generates the example code for a given example
// given a package containing the example code
func ExampleCode(pkg *packages.Package) func(*doc.Example) string {

	return func(ex *doc.Example) string {
		var buf = bytes.NewBuffer(nil)
		buf.WriteString("<details>\n")
		buf.WriteString(fmt.Sprintf("<summary>Example%s</summary>\n\n", ex.Name))
		buf.WriteString(fmt.Sprintf("```go\nfunc Example%s", ex.Name))
		format.Node(buf, pkg.Fset, ex.Code)
		output_lines := strings.Split(ex.Output, "\n")
		buf.WriteString("\n // Output:")
		for _, line := range output_lines {
			buf.WriteString(fmt.Sprintf("\n // %s", line))
		}
		buf.WriteString("\n```\n")
		buf.WriteString("</details>\n")
		return buf.String()
	}
}

// FuncLocation returns the location of the function in a package containing the function
func FuncLocation(pkg *packages.Package) func(*doc.Func) string {

	return func(fn *doc.Func) string {
		var buf = bytes.NewBuffer(nil)
		file := pkg.Fset.File(fn.Decl.Pos())
		start_ln := file.Line(fn.Decl.Pos())
		end_ln := file.Line(fn.Decl.Pos())
		buf.WriteString(fmt.Sprintf("./%s#L%d-L%d", path.Base(file.Name()), start_ln, end_ln))

		return buf.String()
	}
}

func TypeLocation(pkg *packages.Package) func(*doc.Type) string {

	return func(fn *doc.Type) string {
		var buf = bytes.NewBuffer(nil)
		file := pkg.Fset.File(fn.Decl.Pos())
		start_ln := file.Line(fn.Decl.Pos())
		end_ln := file.Line(fn.Decl.Pos())
		buf.WriteString(fmt.Sprintf("./%s#L%d-L%d", path.Base(file.Name()), start_ln, end_ln))

		return buf.String()
	}
}

// FuncSignature returns a function that generates the function signature for a given function in a package
// This function is provided the template parser as 'signature'
// Usage:
// ```cheetah
// //in a template file
// {{signature .Func}} // where .Func is a type of *doc.Func
// ```
// Ex: for this function it would return:
// ```go
// func FuncSignature(pkg *packages.Package) func(*doc.Func) string {
// ```
func FuncSignature(pkg *packages.Package) func(*doc.Func) string {

	return func(fn *doc.Func) string {
		var buf = bytes.NewBuffer(nil)

		buf.WriteString("```go\n")
		if fn.Decl != nil {
			format.Node(buf, pkg.Fset, fn.Decl)
		}
		buf.WriteString("\n```")
		return buf.String()
	}
}

func NoteSignature(pkg *packages.Package) func(string, map[string][]*doc.Note) string {

	return func(name string, notes map[string][]*doc.Note) string {
		_notes, found := notes["NOTE"]
		if !found {
			return ""
		}
		var buf = bytes.NewBuffer(nil)

		var header_written bool
		for _, note := range _notes {

			if note.UID == name {
				if !header_written {
					buf.WriteString(">[!NOTE]\n")
					header_written = true
				}
				buf.WriteString(fmt.Sprintf(">%s", note.Body))
			}
		}

		return buf.String()
	}
}

func TypeSignature(pkg *packages.Package) func(*doc.Type) string {

	return func(fn *doc.Type) string {
		var buf = bytes.NewBuffer(nil)

		buf.WriteString("```go\n")
		if fn.Decl != nil {
			format.Node(buf, pkg.Fset, fn.Decl)
		}
		buf.WriteString("\n```")
		return buf.String()
	}
}

// Generate generates the README.md file for the packages that are loaded
// The README.md file is generated in the directory of the package
// The README.md file is generated using the template file provided or the default template in none is provided
func (readme *Readme) Generate() (err error) {

	for _, pkg := range readme.RefinedPkgs {
		package_readme := PackageReadme{
			P:       pkg,
			Options: *readme.options,
		}
		if package_readme.Pkg, err = doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath); err != nil {
			return
		}
		var readme_file *os.File
		if pkg.GoFiles == nil || len(pkg.GoFiles) == 0 {
			continue
		}
		var readme_file_path = filepath.Dir(pkg.GoFiles[0])
		if readme_file, err = os.Create(path.Join(readme_file_path, "README.md")); err != nil { //fmt.Sprintf("./README.%s.md", pkg.Name)
			return
		}

		fn_map := template.FuncMap{
			"example":        ExampleCode(pkg),
			"fn_signature":   FuncSignature(pkg),
			"fn_location":    FuncLocation(pkg),
			"type_signature": TypeSignature(pkg),
			"type_location":  TypeLocation(pkg),
			"note":           NoteSignature(pkg),
		}
		var tmpl *template.Template
		if readme.options.TemplateFile != "" {
			var template_data []byte
			template_file, err = os.Open(readme.options.TemplateFile)
			if err != nil {
				return err
			}
			template_data, err = io.ReadAll(template_file)
			if err != nil {
				return err
			}

			if tmpl, err = template.New(pkg.Name).Funcs(fn_map).Parse(string(template_data)); err != nil {
				return
			}
		} else {
			if tmpl, err = template.New(pkg.Name).Funcs(fn_map).Parse(readme_template); err != nil {
				return
			}
		}

		if err = tmpl.Execute(readme_file, package_readme); err != nil {
			return
		}
	}
	return
}
