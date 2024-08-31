# godoc_readme
<!-- THIS FILE IS GENERATED. DO NOT EDIT! -->


Godoc-readme is a CLI that generates a README.md file for your go project using comments you already write for godoc!

In fact, this README.md file was generated using godoc-readme! :open_mouth:

> [!Note]
> Adding a `//go:generate godoc-readme directive` to your go file will generate a README.md file for your package when the `go generate` command is run.

Usage:

	godoc-readme [flags]

Flags:

	-h, --help   help for auto-readme
	-r, --recursive   recursively search for go packages in the directory and generate a README.md for each package

more details about the package

## Functions

### [func ExampleCode](./readme.go#L168-L168)
```go
func ExampleCode(pkg *packages.Package) func(*doc.Example) string
```

> ExampleCode returns a function that generates the example code for a given example
given a package containing the example code



### [func Execute](./readme.go#L65-L65)
```go
func Execute(args ...string)
```

> Execute runs the root command using the os.Args by default
Optionally, you can pass in a list of arguments to run the command with



### [func FuncLocation](./readme.go#L187-L187)
```go
func FuncLocation(pkg *packages.Package) func(*doc.Func) string
```

> FuncLocation returns the location of the function in a package containing the function



### [func FuncSignature](./readme.go#L211-L211)
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



## Examples


```go
func Example_godoc_readme{
	Execute("-h")

}
 // Output:
 // 
 // Generate README.md file for your go project using comments you already write for godoc
 // 
 // Usage:
 //   godoc-reademe [flags]
 // 
 // Flags:
 //   -h, --help   help for godoc-reademe
 // 
```



