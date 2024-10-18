package template_functions

import (
	"fmt"
	"go/doc"

	"golang.org/x/tools/go/packages"
)

func Title(pkg *packages.Package, doc *doc.Package) func() string {
	
	return func() string {
		if pkg == nil {
			return fmt.Sprintf("Package `%s`", doc.Name)
		} else {
			// title_pattern := regexp.MustCompile("^@Title{(.*)}$")
			// matches := title_pattern.FindAllString(doc.Doc, -1) 
			// fmt.Println(matches)
			// if len(matches) > 0 {
			// 	return matches[0]
			// }
			title := make([]byte, 0, 256)
			for _, c := range doc.Doc {
				if c == '\n' {
					break
				}
				title = append(title, byte(c))
			}
			if len(title) == 0 {
				return fmt.Sprintf("Package `%s`", doc.Name)
			}
			return string(title)
		}
	}
}