// more details about the package
package godoc_readme

//go:generate go run ./cmd/main.go
import (
	"bytes"
	_ "embed"
	"fmt"
	"go/doc"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "godoc-reademe",
	Short: "Generate README.md file for your go project using comments you already write for godoc",
	Long:  `Generate README.md file for your go project using comments you already write for godoc`,
	Run: func(cmd *cobra.Command, args []string) {
		if readme, err := NewReadme(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			if err = readme.Generate(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
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

//go:embed README.tmpl
var readme_template string

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// Readme is a struct that holds the packages, ast and docs of the package
// And is used to pass data to the readme template
type Readme struct {
	pkgs        []*packages.Package
	RefinedPkgs map[string]*packages.Package
	Docs        []*doc.Package
	options     *ReadmeOptions
}

type ReadmeOptions struct {
	Dir               string `env:"AUTO_README_MODULE_ROOT" default:"."`
	package_load_mode packages.LoadMode
}

func NewReadme(opts ...func(*ReadmeOptions)) (readme *Readme, err error) {
	readme = &Readme{
		options: &ReadmeOptions{
			package_load_mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypes | packages.NeedEmbedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedDeps,
		},
		RefinedPkgs: map[string]*packages.Package{},
	}
	for _, opt := range opts {
		opt(readme.options)
	}
	cfg := &packages.Config{
		Mode:  readme.options.package_load_mode,
		Dir:   readme.options.Dir,
		Tests: true,
	}
	readme.pkgs, err = packages.Load(cfg, "./...")
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

type PackageReadme struct {
	Options ReadmeOptions
	P       *packages.Package
	Pkg     *doc.Package
}

// ExampleCode returns a function that generates the example code for a given example
// given a package containing the example code
func ExampleCode(pkg *packages.Package) func(*doc.Example) string {

	return func(ex *doc.Example) string {
		var buf = bytes.NewBuffer([]byte(fmt.Sprintf("```go\nfunc Example%s", ex.Name)))
		if ex.Code != nil {
			format.Node(buf, pkg.Fset, ex.Code)
		}
		output_lines := strings.Split(ex.Output, "\n")
		buf.WriteString("\n // Output:")
		for _, line := range output_lines {
			buf.WriteString(fmt.Sprintf("\n // %s", line))
		}

		buf.WriteString("\n```")
		return buf.String()
	}
}

// FuncLocation returns the location of the function in a package containing the function
func FuncLocation(pkg *packages.Package) func(*doc.Func) string {

	return func(fn *doc.Func) string {
		var buf = bytes.NewBuffer(nil)
		file := pkg.Fset.File(fn.Decl.Pos())
		start_ln := file.Line(fn.Decl.Type.Pos())
		end_ln := file.Line(fn.Decl.Type.End())
		buf.WriteString(path.Join(pkg.PkgPath, "blob/main", path.Base(file.Name()), fmt.Sprintf("#L%d-L%d", start_ln, end_ln)))

		return buf.String()
	}
}

// FuncSignature returns a function that generates the function signature for a given function in a package
// This function is provided the template parser as 'signature'
// Usage:
// ```go
// //in a template file
// {{signature .Func}} // where .Func is a type of *doc.Func
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
			"example":   ExampleCode(pkg),
			"signature": FuncSignature(pkg),
			"location":  FuncLocation(pkg),
		}
		var tmpl *template.Template
		if tmpl, err = template.New(pkg.Name).Funcs(fn_map).Parse(readme_template); err != nil {
			return
		}

		if err = tmpl.Execute(readme_file, package_readme); err != nil {
			return
		}
	}
	return
}
