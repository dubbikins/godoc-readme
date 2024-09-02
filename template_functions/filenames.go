package template_functions

import (
	"os"
	"strings"
)

func RelativeFilename(abs string) (relative string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Replace(abs, pwd, ".", 1)
}
