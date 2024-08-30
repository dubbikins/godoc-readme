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

### [func ExampleCode](https://github.com/dubbikins/godoc-readme/blob/main/readme.go/#L135-L135)
```go
func ExampleCode(pkg *packages.Package) func(*doc.Example) string
```

ExampleCode returns a function that generates the example code for a given example
given a package containing the example code



### [func Execute](https://github.com/dubbikins/godoc-readme/blob/main/readme.go/#L41-L41)
```go
func Execute(args ...string)
```

Execute runs the root command using the os.Args by default
Optionally, you can pass in a list of arguments to run the command with



### [func FuncLocation](https://github.com/dubbikins/godoc-readme/blob/main/readme.go/#L154-L154)
```go
func FuncLocation(pkg *packages.Package) func(*doc.Func) string
```

FuncLocation returns the location of the function in a package containing the function



### [func FuncSignature](https://github.com/dubbikins/godoc-readme/blob/main/readme.go/#L174-L174)
```go
func FuncSignature(pkg *packages.Package) func(*doc.Func) string
```

FuncSignature returns a function that generates the function signature for a given function in a package
This function is provided the template parser as 'signature'
Usage:
```go
//in a template file
{{signature .Func}} // where .Func is a type of *doc.Func
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



