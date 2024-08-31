# Package `godoc_readme`
<!-- THIS FILE IS GENERATED. DO NOT EDIT! -->


Godoc-readme is a CLI that generates a README.md file for your go project using comments you already write for godoc!

In fact, this README.md file was generated using godoc-readme! :open_mouth:

> [!Note]
> Adding a `//go:generate godoc-readme directive` to your go file will generate a README.md file for your package when the `go generate` command is run.

Usage:

	godoc-readme [flags]

Flags:

	-h, --help              help for godoc-readme
	-r, --recursive         Recursively search for go packages in the directory and generate a README.md for each package (default true)
	-t, --template string   The template file to use for generating the README.md file

more details about the package


## Types

### [type PackageReadme](./readme.go#L161-L161)
```go
type PackageReadme struct {
	Options ReadmeOptions
	P       *packages.Package
	Pkg     *doc.Package
}
```
> PackageReadme is a struct that holds the package, ast and docs of the package
It's used to pass data to the readme template




### [type Readme](./readme.go#L78-L78)
```go
type Readme struct {
	RefinedPkgs map[string]*packages.Package
	Docs        []*doc.Package
	// contains filtered or unexported fields
}
```
> Readme is a struct that holds the packages, ast and docs of the package
And is used to pass data to the readme template


#### Functions

### [func NewReadme](./readme.go#L107-L107)
```go
func NewReadme(opts ...func(*ReadmeOptions)) (readme *Readme, err error)
```



### [type ReadmeOptions](./readme.go#L87-L87)
```go
type ReadmeOptions struct {
	Dir          string `env:"GODOC_README_MODULE_DIR"`
	DirPattern   string `env:"GODOC_README_MODULE_DIR_PATTERN" default:"./..."`
	TemplateFile string `env:"GODOC_README_TEMPLATE_FILE"`
	// contains filtered or unexported fields
}
```
> ReadmeOptions is a struct that holds the options for the Readme struct
You can set the options via the options functions or by setting the environment variables defined in the `env` struct tag for the Option field






## Functions

### [func ExampleCode](./readme.go#L169-L169)
```go
func ExampleCode(pkg *packages.Package) func(*doc.Example) string
```

> ExampleCode returns a function that generates the example code for a given example
given a package containing the example code



### [func Execute](./readme.go#L66-L66)
```go
func Execute(args ...string)
```

> Execute runs the root command using the os.Args by default
Optionally, you can pass in a list of arguments to run the command with



### [func FuncLocation](./readme.go#L189-L189)
```go
func FuncLocation(pkg *packages.Package) func(*doc.Func) string
```

> FuncLocation returns the location of the function in a package containing the function



### [func FuncSignature](./readme.go#L226-L226)
```go
func FuncSignature(pkg *packages.Package) func(*doc.Func) string
```

> FuncSignature returns a function that generates the function signature for a given function in a package
This function is provided the template parser as 'signature'
Usage:
```cheetah
//in a template file
{{signature .Func}} // where .Func is a type of *doc.Func
```
Ex: for this function it would return:
```go
func FuncSignature(pkg *packages.Package) func(*doc.Func) string {
```



### [func TypeLocation](./readme.go#L202-L202)
```go
func TypeLocation(pkg *packages.Package) func(*doc.Type) string
```

> 


### [func TypeSignature](./readme.go#L240-L240)
```go
func TypeSignature(pkg *packages.Package) func(*doc.Type) string
```

> 



## Examples

<details>
<summary>Example</summary>

```go
func Example{
	Execute()

}
 // Output:
 // 
 // README.md file generated successfully :tada:
 // 
```
</details>


<details>
<summary>Example_help_command</summary>

```go
func Example_help_command{
	Execute("-h")

}
 // Output:
 // 
 // Generate README.md file for your go project using comments you already write for godoc
 // 
 // Usage:
 //   godoc-readme [flags]
 // 
 // Flags:
 //   -h, --help   help for godoc-readme
 //   -r, --recursive   Recursively search for go packages in the directory and generate a README.md for each package (default true)
 //   -t, --template string   The template file to use for generating the README.md file
 // 
```
</details>


<details>
<summary>Example_template_file</summary>

```go
func Example_template_file{
	Execute("-t", "README.tmpl")

}
 // Output:
 // 
 // README.md file generated successfully :tada:
 // 
```
</details>






