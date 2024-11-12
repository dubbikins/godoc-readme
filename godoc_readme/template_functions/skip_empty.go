package template_functions


func SkipEmpty(skip_empty bool) func(string ) bool {
	
	return func(doc_string string) bool {
		return skip_empty  && doc_string == ""
	}
}