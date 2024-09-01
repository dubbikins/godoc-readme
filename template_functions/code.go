package template_functions

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

/*
CodeBlock returns the start (or end) of a code block in markdown
If you provide a language, it will be used to specify the language format of the code block
In a template, you should call this function like this: `{{ code_block "go" }}` or `{{ code_block }}` to start the code block
and `{{ code_block }}` to end the code block after you've rendered the contens between the start and end of the code block

Example:

```tmpl
{{ code_block "go" }}
package main

	func main() {
		fmt.Println("Hello, World!")
	}

{{ code_block }}
```
This will render the following markdown:

```go
package main

	func main() {
		fmt.Println("Hello, World!")
	}

```
*/
func CodeBlock(pkg *packages.Package) func(lang ...string) string {
	return func(lang ...string) string {
		var _lang string
		if l := len(lang); l > 0 {
			_lang = lang[0]
		}
		return fmt.Sprintf("```%s", _lang)
	}
}
