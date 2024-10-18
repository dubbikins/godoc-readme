package godoc_readme

import (
	"bytes"
	"embed"
	"fmt"
	"go/doc"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/dubbikins/envy"
	"github.com/dubbikins/godoc-readme/godoc_readme/template_functions"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/tools/go/packages"
)

// The readme templates are embedded in the binary so that it can be used as a default template
// This value can be overridden by providing a template file using the --template flag or the GODOC_README_TEMPLATE_FILE environment variable
//
//go:embed templates/*
var readme_templates embed.FS


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
	TestPkgs    map[string]*packages.Package
	pkgs    []*packages.Package
	options *ReadmeOptions
	stdio io.Reader
	readmes []*PackageReadme
	// buf bytes.Buffer
}

type ReadmeOptions struct {

	Dir				   string`env:"GODOC_README_MODULE_DIR"`
	DirPattern         string `env:"GODOC_README_MODULE_DIR_PATTERN" default:"./..."`
	TemplateFile       string `env:"GODOC_README_TEMPLATE_FILE"`
	Format             func([]byte) []byte
	package_load_mode  packages.LoadMode
	
	ConfirmUpdates bool
	Flags template_functions.Flags
}



// ReadmeOptions is a struct that holds the options for the Readme struct
// You can set the options via the options functions or by setting the environment variables defined in the `env` struct tag for the Option field


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
		TestPkgs: map[string]*packages.Package{},
		readmes: []*PackageReadme{},
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
		for f_index, file := range pkg.GoFiles {
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
		// if p, found := readme.Pkgs[pkg.Name]; found {
			
		// }
		if !strings.Contains(pkg.ID, "test") {
			readme.TestPkgs[pkg.Name] = pkg
		}else {
			readme.Pkgs[pkg.Name] = pkg
		}	
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
	bytes.Buffer
	package_dir string
	file_path string
	rel_file_path string
	file_name string
	file *os.File
	cwd string
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
	fmt.Println("Generated Documentation:")
	for _, pkg := range readme.Packages {
		var pkg_readme *PackageReadme
		if pkg_readme, err = readme.generate_pkg_readme(pkg,"README.md"); err != nil {
			return
		}
		readme.readmes = append(readme.readmes, pkg_readme)
		fmt.Printf("\t- ./%s\n", pkg_readme.rel_file_path)
	}
	return
}

func (readme *Readme) PackageREADMES(yield func(*PackageReadme) bool) {
	for _, _readme := range readme.readmes {
		if !yield(_readme) {
			break
		}
	}
}

func (readme *Readme) Packages(yield func(string, *packages.Package) bool) {
	var pkgs map[string] *packages.Package
	if readme.options.Flags.SkipExamples {
		pkgs = readme.Pkgs
	}else {
		pkgs =readme.TestPkgs
	}
	var sorted_pkg_keys []string  = make([]string,0, len(pkgs))
	for name, _ := range pkgs {
		sorted_pkg_keys = append(sorted_pkg_keys, name)
	}
	sort.Slice(sorted_pkg_keys, func(i, j int) bool {
		return strings.Compare(sorted_pkg_keys[i], sorted_pkg_keys[j])  == -1
	})
	for _, key := range sorted_pkg_keys {
		if !yield(key, pkgs[key]) {
			break
		}
	}
}

func (readme *Readme) generate_pkg_readme(pkg *packages.Package, filename string) (package_readme *PackageReadme, err error) {
		package_readme = &PackageReadme{
			Pkg:     pkg,
			Options: *readme.options,
		}
		if package_readme.Doc, err = doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath,doc.AllDecls| doc.AllMethods |doc.PreserveAST); err != nil {
			return
		}
		// var buf = bytes.NewBuffer(nil)	
		if len(pkg.GoFiles) == 0 {
			return
		}
		var readme_file_path = filepath.Dir(pkg.GoFiles[0])
		package_readme.file_name = path.Join(readme_file_path,filename )
		var existing_file *os.File
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
			"title":         template_functions.Title(pkg, package_readme.Doc),
			"flags":         template_functions.GetFlag(readme.options.Flags),
			"filename":          filepath.Base,
		}
		var tmpl *template.Template
		// if readme.options.TemplateFile != "" {
			// var template_data []byte
			// template_file, err = os.Open(readme.options.TemplateFile)
			// if err != nil {
			// 	return 
			// }
			// template_data, err = io.ReadAll(template_file)
			// if err != nil {
			// 	return 
			// }
			// if tmpl, err = template.New(pkg.Name).Funcs(fn_map).Parse(string(template_data)); err != nil {
			// 	return
			// }
		//} else {
			if tmpl, err = template.New("README.tmpl").Funcs(fn_map).ParseFS(readme_templates, "templates/*.tmpl"); err != nil {
				return
			}
			
		//}
		if err = tmpl.Execute(package_readme, package_readme); err != nil {
			return
		}
		if readme.options.Format == nil {
			readme.options.Format = FormatMarkdown
		}
		if _, err = os.Stat(package_readme.file_name); err == nil {
			if readme.options.ConfirmUpdates {
				if existing_file, err = os.Open(package_readme.file_name); err != nil {
					return
				}
				defer existing_file.Close()
				fmt.Printf("The file %s already exists, do you want to overwrite it? (y/n/diff): ", package_readme.file_name)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) == "y" {
					fmt.Println("Overwriting file...")
				} else if strings.ToLower(response) == "diff" {
					fmt.Println("Generating file diff...")
					dmp := diffmatchpatch.New()
					var existing_data []byte 
					if existing_data, err = io.ReadAll(existing_file); err != nil {
						return
					}
					diffs := dmp.DiffMain(string(existing_data), package_readme.String(), false)
					fmt.Println("Original:")
					fmt.Println(string(existing_data))
					fmt.Println("New:")
					fmt.Println(dmp.DiffPrettyText(diffs))
					fmt.Println("Proceed with overwritting file? (y/n): ", package_readme.file_name)
					var response string
					fmt.Scanln(&response)
					if strings.ToLower(response) != "y" {
						return
					} 
				} else {
					fmt.Println("Not Overwriting file...")
					return
				}
			}
		}
		if package_readme.file, err = os.Create(package_readme.file_name); err != nil { //fmt.Sprintf("./README.%s.md", pkg.Name)
			return
		}
		//defer package_readme.file.Close()
		if _, err = package_readme.file.Write(readme.options.Format(package_readme.Bytes())); err != nil {
			return
		}
		if  package_readme.cwd , err = os.Getwd(); err != nil {
			return
		}
		// fmt.Println("CWD: ", package_readme.cwd)
		// fmt.Println("Module Dir: ", package_readme.Pkg.Module.Dir)
		// fmt.Println("File Name: ", package_readme.file_name)
		// fmt.Println("PkgPath: ", package_readme.Pkg.PkgPath)
		// fmt.Println("Module Path: ", package_readme.Pkg.Module.Path)

		// if  package_readme.rel_file_path, err = filepath.Rel( package_readme.cwd,package_readme.Pkg.Module.Dir);err != nil {
		// 	return
		// }
		package_readme.rel_file_path = filepath.Join(strings.Replace( package_readme.Pkg.PkgPath, package_readme.Pkg.Module.Path, "./", 1), filename)
		return 
}