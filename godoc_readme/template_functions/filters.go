package template_functions

import "go/doc"

type MethodsOptions struct {
	SkipEmpty bool
}

func FilteredFuncs(options MethodsOptions) func(funcs[]*doc.Func) []*doc.Func {
	
	return func(funcs[]*doc.Func) (filtered_funcs []*doc.Func ){
		filtered_funcs = make([]*doc.Func, 0, len(funcs))
		for _, _func := range funcs {
			if !options.SkipEmpty || _func.Doc != "" {
				filtered_funcs = append(filtered_funcs, _func)
			}
		}
		return
	}
}