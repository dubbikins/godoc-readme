package godoc_readme

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"go/doc"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/dubbikins/envy"
	"github.com/dubbikins/godoc-readme/godoc_readme/template_functions"
	"github.com/pkg/browser"
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
	readmes []*PackageReadme
	confirmation_listener net.Listener
	confirmation_listener_port int
	confirmation_server *http.Server
}

// ReadmeOptions is a struct that holds the options for the Readme struct
// You can set the options via the options functions or by setting the environment variables defined in the `env` struct tag for the Option field
type ReadmeOptions struct {
	PackageDir string 
	Dir				   string
	Format             func([]byte) []byte
	package_load_mode  packages.LoadMode
	Env  []string `env:"-"`
	ConfirmUpdates bool
	Flags template_functions.Flags
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
		Env:  append(os.Environ(), readme.options.Env...),
		Tests: true,
	}, readme.options.PackageDir); err != nil {
		return
	}
	if packages.PrintErrors(readme.pkgs) > 0 {
		// Package failed to parse
		os.Exit(1)
	}
	if readme.options.ConfirmUpdates {
		readme.confirmation_listener_port = 8080
		readme.confirmation_server = &http.Server{}
		if readme.confirmation_listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", readme.confirmation_listener_port)); err != nil {
			panic(err)
		}	
		
		
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
	
		if !strings.Contains(pkg.ID, "test") {
			if _, exists := readme.Pkgs[pkg.Name]; exists {
				continue
			}
			readme.Pkgs[pkg.Name] = pkg	
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
	rel_file_path string
	file_name string
	file *os.File
	cwd string
	rejected bool
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
	if readme.options.ConfirmUpdates {
		go func() {
			err = readme.confirmation_server.Serve(tcpKeepAliveListener{readme.confirmation_listener.(*net.TCPListener)})
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}

		}()
	}
	for _, pkg := range readme.Packages {
		var pkg_readme *PackageReadme
		if pkg_readme, err = readme.generate_pkg_readme(pkg,"README.md"); err != nil {
			return
		}
		readme.readmes = append(readme.readmes, pkg_readme)
	}
	fmt.Println("Results:")
	for _readme := range readme.READMES {
		if !_readme.rejected  || !readme.options.ConfirmUpdates {
			fmt.Printf("\t- %q \u2705\n", _readme.file_name)
		}else {
			fmt.Printf("\t- %q \u274c\n", _readme.file_name)
		}
	
	}
	return
}

func (readme *Readme) READMES(yield func(*PackageReadme) bool) {
	for _, _readme := range readme.readmes {
		if !yield(_readme) {
			break
		}
	}
}

func (readme *Readme) Packages(yield func(string, *packages.Package) bool) {
	var sorted_pkg_keys []*packages.Package = make([]*packages.Package,0, len(readme.Pkgs))
	for _, pkg := range readme.Pkgs {
		sorted_pkg_keys = append(sorted_pkg_keys, pkg)
	}
	sort.Slice(sorted_pkg_keys, func(i, j int) bool {
		return strings.Compare(sorted_pkg_keys[i].Name, sorted_pkg_keys[j].Name)  == -1
	})
	for _, pkg := range sorted_pkg_keys {
		if !yield(pkg.Name, pkg) {
			break
		}
	}
}
func (readme *Readme) template_functions (package_readme *PackageReadme) template.FuncMap {
	return template.FuncMap{
		"example":       template_functions.ExampleCode(package_readme.Pkg),
		"skip_empty":    template_functions.SkipEmpty(readme.options.Flags.SkipEmpty),
		"filtered_funcs":    template_functions.FilteredFuncs(template_functions.MethodsOptions{
			SkipEmpty: readme.options.Flags.SkipEmpty,
		}),
		"code":          template_functions.CodeBlock(package_readme.Pkg),
		"fmt":           template_functions.FormatNode(package_readme.Pkg),
		"link":          template_functions.Link(package_readme.Pkg),
		"alert":         template_functions.Alert(package_readme.Pkg, package_readme.Doc.Notes),
		"doc":           template_functions.DocString,
		"gen_decl": 	 template_functions.GenDeclaration(package_readme.Pkg),
		"spec_decl": 	 template_functions.SpecDeclaration(package_readme.Pkg),
		"fn_decl": 		 template_functions.FuncDeclaration(package_readme.Pkg),
		"decl":          template_functions.Declaration(package_readme.Pkg),
		"section":       template_functions.Section,
		"pkg_doc":       template_functions.PackageDocString,
		"relative_path": template_functions.RelativeFilename,
		"title":         template_functions.Title(package_readme.Pkg, package_readme.Doc),
		"flags":         template_functions.GetFlag(readme.options.Flags),
		"filename":          filepath.Base,
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
		if len(pkg.GoFiles) == 0 {
			return
		}
		var readme_file_path = filepath.Dir(pkg.GoFiles[0])
		package_readme.file_name = path.Join(readme_file_path,filename )
		var tmpl *template.Template
		if tmpl, err = template.New("README.tmpl").Funcs(readme.template_functions(package_readme)).ParseFS(readme_templates, "templates/*.tmpl"); err != nil {
			return
		}	
		if err = tmpl.Execute(package_readme, package_readme); err != nil {
			return
		}
		if readme.options.Format == nil {
			readme.options.Format = FormatMarkdown
		}

		if !readme.confirm_changes(package_readme) {
			package_readme.rejected = true
			return
		}
		if package_readme.file, err = os.Create(package_readme.file_name); err != nil { 
			return
		}
		if _, err = package_readme.file.Write(readme.options.Format(package_readme.Bytes())); err != nil {
			return
		}
		if  package_readme.cwd , err = os.Getwd(); err != nil {
			return
		}

		package_readme.rel_file_path = filepath.Join(strings.Replace( package_readme.Pkg.PkgPath, package_readme.Pkg.Module.Path, "./", 1), filename)
		return 
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (readme *Readme) confirm_changes(package_readme *PackageReadme) (confirm bool) {
	var existing_file *os.File
	if !readme.options.ConfirmUpdates {
		return true
	}
	if _, err := os.Stat(package_readme.file_name); err == nil {
		if readme.options.ConfirmUpdates {
			if existing_file, err = os.Open(package_readme.file_name); err != nil {
				return
			}
			defer existing_file.Close()
			dmp := diffmatchpatch.New()
			var existing_data []byte 
			if existing_data, err = io.ReadAll(existing_file); err != nil {
				return
			}
			diffs := dmp.DiffMain(string(existing_data), string(readme.options.Format(package_readme.Bytes())), false)
			// var differ *bytes.Buffer
			confirmation_response := make(chan bool, 1)
			readme.confirmation_server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {

			case "/":
				tmpl, err := template.New(".diff.template").Funcs(template.FuncMap{
					"hostname": func() string {
						return readme.confirmation_listener.Addr().String()
					},
				}).ParseFS(readme_templates, "templates/.diff.template")
				if err != nil {
					fmt.Println(err)
					return
				}
				w.Header().Add("Content-Type", "text/html")
				if err = tmpl.Execute(w, diffs); err != nil {
					fmt.Println(err)
					return
				}
				return
			case "/confirm":
				fmt.Println("Confirmation signal received from the browser")
				w.Write([]byte(fmt.Sprintf(`<html><head><script>window.onload = ()=>{var self = top.window.open('', '_self', ''); window.location.href="vscode://file%s"; setTimeout(()=>  self.close(),2000);}</script></head><body>Confirmed. Redirecting to editor... Window will close shortly</html>`, package_readme.file_name)))
				confirmation_response<- true
			case "/reject":
				fmt.Println("Rejection signal received from the browser")
				w.Write([]byte(fmt.Sprintf(`<html><head><script>window.onload = ()=>{var self = top.window.open('', '_self', ''); window.location.href="vscode://file%s"; setTimeout(()=>  self.close(),2000);}</script></head><body>Rejected. Redirecting to editor... Window will close shortly</html>`, package_readme.file_name)))
				confirmation_response<- false
			}
			})
			go func () {
				var response string
				fmt.Print("Proceed with overwritting README? (y/n): ", package_readme.file_name)
				fmt.Scanln(&response)
				if strings.ToLower(response) == "y" {
					confirmation_response<- true
				}else {
					confirmation_response<- false
				}
			}()
			<-time.After(1 * time.Second)
			fmt.Println("Viewing changes to " + package_readme.file_name + " in the browser...")
			if err = browser.OpenURL(fmt.Sprintf("http://%s", readme.confirmation_listener.Addr().String())); err != nil {
				fmt.Println(err)
				return
			}
			select {
			case confirm = <-confirmation_response:
			case <-time.After(5 * time.Minute):
				fmt.Println("Confirmation timed out")
			}
			if !confirm {
				fmt.Printf("Changes to %q rejected \u274c\n", package_readme.file_name)
				return
			}
			fmt.Printf("Changes to %q accepted \u2705\n", package_readme.file_name)
		}
	}
	return
}