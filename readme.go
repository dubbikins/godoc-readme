package godoc_readme

import (
	"embed"
	"fmt"
	"go/doc"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dubbikins/envy"
	"github.com/dubbikins/godoc-readme/template_functions"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

// The readme templates are embedded in the binary so that it can be used as a default template
// This value can be overridden by providing a template file using the --template flag or the GODOC_README_TEMPLATE_FILE environment variable
//
//go:embed templates/*
var readme_templates embed.FS
var recursive bool = true
var template_filename string
var template_file *os.File

//go:generate go run ./cmd/main.go

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

// Readme is a struct that holds the packages, ast and docs of the package
// And is used to pass data to the readme template
/*
```mermaid
classDiagram
	class Readme
	Readme : +map[string]*packages.Package Pkgs
    Readme : -[]*packages.Package pkgs
	Readme : -ReadmeOptions options
	Readme : +Generate() error
	Readme --> ReadmeOptions
	Readme --> PackageReadme
	class ReadmeOptions
	ReadmeOptions : -string Dir
	ReadmeOptions : -string DirPattern
	ReadmeOptions : -string TemplateFile
	class PackageReadme
	PackageReadme : +ReadmeOptions Options
	PackageReadme : +packages.Package Pkg
	PackageReadme : +doc.Package Doc
```
*/
type Readme struct {
	Pkgs    map[string]*packages.Package
	pkgs    []*packages.Package
	options *ReadmeOptions
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
		Pkgs: map[string]*packages.Package{},
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
		if p, found := readme.Pkgs[pkg.Name]; found {
			if !strings.Contains(p.ID, "test") {
				readme.Pkgs[pkg.Name] = pkg
			}
		}
		readme.Pkgs[pkg.Name] = pkg
	}
	return
}

// PackageReadme is a struct that holds the package, ast and docs of the package
// It's used to pass data to the readme template
type PackageReadme struct {
	Options ReadmeOptions
	Pkg     *packages.Package
	Doc     *doc.Package
}

/*
Generate creates the README.md file for the packages that are registered with a `Readme`

The README is generated in the directory of the package using the template file provided or the default template in none is provided.
The template functions that are made available to the template arg defined in the [`template_functions` package](./template_functions/README.go)
*/
func (readme *Readme) Generate() (err error) {

	for _, pkg := range readme.Pkgs {
		package_readme := PackageReadme{
			Pkg:     pkg,
			Options: *readme.options,
		}
		if package_readme.Doc, err = doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath); err != nil {
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
			"example": template_functions.ExampleCode(pkg),
			"code":    template_functions.CodeBlock(pkg),
			"fmt":     template_functions.FormatNode(pkg),
			"link":    template_functions.Link(pkg),
			"Alert":   template_functions.Alert(pkg, package_readme.Doc.Notes),
			"section": template_functions.Section,
			"doc":     template_functions.DocString,
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
			if tmpl, err = template.New(".godoc-readme.tmpl").Funcs(fn_map).ParseFS(readme_templates, "templates/*.tmpl"); err != nil {
				return
			}
		}

		if err = tmpl.Execute(readme_file, package_readme); err != nil {
			return
		}
	}
	return
}
