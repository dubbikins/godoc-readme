package template_functions

import "strings"

// Section returns a copy of doc with `n` number of `>`'s added to the beginning of each line.
// You can call this function in a template by using `{{ section .Doc 1 }}` where `.Doc` is *string* field.
// Example:
//
//	`Section("This is a section", 1)` returns "> This is a section"
func Section(doc string, n int) string {
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(">", n) + line
	}
	return strings.Join(lines, "\n")
}
