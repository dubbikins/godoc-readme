/*
Auto-Readme is a CLI that generates a README.md file for your go project using comments you already write for godoc!

Usage:

	auto-readme [flags]

Flags:

	-h, --help   help for auto-readme
	-r, --recursive   recursively search for go packages in the directory and generate a README.md for each package

> [!Note]
> Adding a `//go:generate auto-readme directive` to your go file will generate a README.md file for your package when the `go generate` command is run.
*/
package autoreadme

//go:generate go run ./cmd/main.go
import (
	_ "embed"
	"fmt"
	"go/doc"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "auto-readme",
	Short: "Generate README.md file for your go project using comments you alreayd write for godoc",
	Long:  `Generate README.md file for your go project using comments you alreayd write for godoc`,
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

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//go:embed README.tmpl
var readme_template string

// Readme is a struct that holds the packages, ast and docs of the package
// And is used to pass data to the readme template
type Readme struct {
	Pkgs    []*packages.Package
	Docs    []*doc.Package
	options *ReadmeOptions
}

type ReadmeOptions struct {
	Dir               string `env:"AUTO_README_MODULE_ROOT" default:"."`
	package_load_mode packages.LoadMode
}

func NewReadme(opts ...func(*ReadmeOptions)) (readme *Readme, err error) {
	readme = &Readme{
		options: &ReadmeOptions{
			package_load_mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedImports | packages.NeedTypes,
		},
	}
	for _, opt := range opts {
		opt(readme.options)
	}
	if err = readme.loadPackages(); err != nil {
		return
	}

	readme.Docs = make([]*doc.Package, len(readme.Pkgs))
	for i, pkg := range readme.Pkgs {
		log.Printf("Package %q\n", pkg.Name)
		log.Println("Go Files:")

		for _, file := range pkg.GoFiles {
			log.Printf("\t- %q\n", file)

		}

		if readme.Docs[i], err = doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath); err != nil {
			return
		}

		log.Println("Go Docs:")
		log.Println(readme.Docs[i].Doc)
		for _, fn := range readme.Docs[i].Funcs {
			log.Printf("\t- %q\n", fn.Doc)
		}

	}

	// for i, file := range readme.Pkg.Syntax {

	// 	readme.AST.Files[readme.Pkg.GoFiles[i]] = file
	// }

	//
	// fmt.Println(readme.Docs.Doc)

	return
}

func (r *Readme) loadPackages() (err error) {

	cfg := &packages.Config{
		Mode: r.options.package_load_mode,
		Dir:  r.options.Dir,
	}
	r.Pkgs, err = packages.Load(cfg, ".")
	if err != nil {
		return fmt.Errorf("failed to load package in dir %q: %w", cfg.Dir, err)
	}
	if packages.PrintErrors(r.Pkgs) > 0 {
		// Package failed to parse
		os.Exit(1)
	}

	return
}

type PackageReadme struct {
	Options ReadmeOptions
	Pkg     *doc.Package
}

func (readme *Readme) Generate() (err error) {
	for i, pkg := range readme.Pkgs {
		var readme_file *os.File
		if readme_file, err = os.Create("./README.md"); err != nil {
			return
		}
		log.Printf("Package %q\n", pkg.Name)
		log.Println("Go Files:")
		var tmpl *template.Template
		if tmpl, err = template.New(pkg.Name).Parse(readme_template); err != nil {
			return
		}
		// var w *bytes.Buffer = bytes.NewBuffer(nil)
		package_readme := PackageReadme{
			Options: *readme.options,
			Pkg:     readme.Docs[i],
		}
		//readme.Docs[i].Types
		//readme.Docs[i].Funcs
		if err = tmpl.Execute(readme_file, package_readme); err != nil {
			return
		}
	}
	return
}
