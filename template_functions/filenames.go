package template_functions

import (
	"fmt"
	"path"

	"golang.org/x/tools/go/packages"
)

func RelativeFilename(pkg *packages.Package) func(abs string) (relative string) {
	return func(abs string) (relative string) {
		
		return fmt.Sprintf("./%s", path.Base(abs))
	}
}
