package template_functions

import (
	"fmt"
	"path"
)

func RelativeFilename(filepath string) (relative string) {
	return fmt.Sprintf("./%s", path.Base(filepath))
}

