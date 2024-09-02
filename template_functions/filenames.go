package template_functions

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func RelativeFilename(pkg *packages.Package) func(abs string) (relative string) {
	return func(abs string) (relative string) {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return strings.Replace(strings.Replace(abs, pwd, ".", 1), fmt.Sprintf("./%s",pkg.PkgPath), "", 1)
	}
}
