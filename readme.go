package godoc_readme

import (
	"bytes"
	"embed"
	"fmt"
	"go/doc"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

//go:generate go run ./godoc-readme/main.go

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
		
	},
}

// Execute runs the root command using the os.Args by default
// Optionally, you can pass in a list of arguments to run the command with
func Execute(args ...string) {
	if len(args) > 0 {
		rootCmd.SetArgs(args)
	}
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

//NOTE(Readme): Because of the simpicity of godoc-readme's templating engine, you can add powerful customizations to your documentation like the class diagram that was created using a code block and the [mermaid.js](https://mermaid.js.org/) library that is supported out of the box with Github markdown. (not all features are supported though.)

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

	Dir				   string`env:"GODOC_README_MODULE_DIR"`
	DirPattern         string `env:"GODOC_README_MODULE_DIR_PATTERN" default:"./..."`
	TemplateFile       string `env:"GODOC_README_TEMPLATE_FILE"`
	Format             func([]byte) []byte
	package_load_mode  packages.LoadMode
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
	//packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypes | packages.NeedEmbedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedDeps
	readme = &Readme{
		options: &ReadmeOptions{
			package_load_mode: ^packages.LoadMode(0),
			Format: FormatMarkdown,
		},
		Pkgs: map[string]*packages.Package{},
	}
	if err = envy.Unmarshal(readme.options); err != nil {
		return
	}
	for _, opt := range opts {
		opt(readme.options)
	}
	if readme.pkgs, err = packages.Load(&packages.Config{
		Mode:  readme.options.package_load_mode,
		Dir:   readme.options.Dir,
		Tests: true,
	}, readme.options.DirPattern); err != nil {
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

// FormatMarkdown applies the following formatting to the markdown:
// 1. Replace all hard-tabs(`\t`) with 4 single space characters (`    `)
// 2. Remove leading whitespace from blank lines
// 3. Replace multiple `\n`(3+) with a single `\n`
func FormatMarkdown(md []byte) []byte {
	md = bytes.ReplaceAll(md, []byte("\t"), []byte("    "))
	var whitespace_blank_line_pattern = regexp.MustCompile(`\n\s*\n`)
	md = whitespace_blank_line_pattern.ReplaceAll(md, []byte("\n\n")) // Remove leading whitespace from blank lines
	var multiple_blank_line_patterns = regexp.MustCompile(`\n\n+\n`)
	md = multiple_blank_line_patterns.ReplaceAll(md, []byte("\n\n"))
	
	return md
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
The following template functions available in the template engine are defined in the [`template_functions` package](./template_functions):
| Function | Description | Example | Output |
| --- | --- | --- | --- |
| `example` | Renders a markdown representation of a `[doc.Example]` instance | `{{ example . }}` where `.` is a [doc.Example]| renders [an example like](/#Examples) |
| `code` | Renders the start (or end) of a code block in markdown, optionally specifying the language format of the code block | `{{ code "go" }}fmt.Println("Hello World"){{ code }}` | `` ```go\nfmt.Println("Hello World")\n```\n`` |
| `fmt` | Renders a formatted string representation of an [ast.Node] | `{{ fmt . }}` | `N/A` |
| `link` | Renders a markdown link to the location of the [ast.Node] in a package | `{{ link "title" . }}` | `[title](...)` where ... is the relative link to the file ,including line numbers |
| `alert` | Renders a markdown alert message based on the notes provided in the [doc.Package] | `{{ alert . "title" }}` | renders the alerts with the "title" target |
| `section` | Renders an indented markdown section header | `{{ section "line 1 text\nline 2 text" 1}}` | `>line 1 text\n>line 2 text` |
| `doc` | Renders a ***package's*** doc string, including in-line alerts, package ref's, etc | `{{ doc . }}` | `N/A` |
| `relative_path` | Replaces the pwd the `.` | `{{ relative_path "/abs/path" }}` where `/abs` is the pwd | returns `./path` |

Additionally, the following functions are available in the template engine:

- `base`: [filepath.Base] Returns the base name of a file path
*/
func (readme *Readme) Generate() (err error) {

	for _, pkg := range readme.Pkgs {
		package_readme := PackageReadme{
			Pkg:     pkg,
			Options: *readme.options,
		}
		if package_readme.Doc, err = doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath,doc.AllDecls| doc.AllMethods |doc.PreserveAST); err != nil {
			return
		}
		var buf = bytes.NewBuffer(nil)
		var readme_file *os.File
		if pkg.GoFiles == nil || len(pkg.GoFiles) == 0 {
			continue
		}
		var readme_file_path = filepath.Dir(pkg.GoFiles[0])
		var file_name = path.Join(readme_file_path, "README.md")
		if readme_file, err = os.Create(file_name); err != nil { //fmt.Sprintf("./README.%s.md", pkg.Name)
			return
		}

		fn_map := template.FuncMap{
			"example":       template_functions.ExampleCode(pkg),
			"code":          template_functions.CodeBlock(pkg),
			"fmt":           template_functions.FormatNode(pkg),
			"link":          template_functions.Link(pkg),
			"alert":         template_functions.Alert(pkg, package_readme.Doc.Notes),
			"doc":           template_functions.DocString,
			"gen_decl": 	 template_functions.GenDeclaration(pkg),
			"spec_decl": 	 template_functions.SpecDeclaration(pkg),
			"fn_decl": 		 template_functions.FuncDeclaration(pkg),
			"decl":          template_functions.Declaration(pkg),
			"section":       template_functions.Section,
			"pkg_doc":       template_functions.PackageDocString,
			"relative_path": template_functions.RelativeFilename,
			"filename":          filepath.Base,
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
			if tmpl, err = template.New("README.tmpl").Funcs(fn_map).ParseFS(readme_templates, "templates/*.tmpl"); err != nil {
				return
			}
		}
		fmt.Println(package_readme.Doc.Consts)
		if err = tmpl.Execute(buf, package_readme); err != nil {
			return
		}
		if readme.options.Format == nil {
			readme.options.Format = FormatMarkdown
		}
		if _, err = readme_file.Write(readme.options.Format(buf.Bytes())); err != nil {
			return
		}
		fmt.Printf("[%s] generated successfully :tada:\n",file_name )
	}
	return
}
