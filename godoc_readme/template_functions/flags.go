package template_functions


type Flags struct {

	SkipImports bool
	SkipExamples bool
	SkipVars bool
	SkipTypes bool
	SkipFuncs bool
	SkipMethods bool
	SkipFilenames bool
	SkipConsts bool
	SkipEmpty bool
	SkipAll bool
}

func GetFlag(flags Flags) func(string ) bool {
	return func(flag_name string) bool {
		switch flag_name {
			case "ShowImports":
				return !flags.SkipImports
			case "ShowExamples":
				return !flags.SkipExamples
			case "ShowVars":
				return !flags.SkipVars
			case "ShowTypes":
				return !flags.SkipTypes
			case "ShowFuncs":
				return !flags.SkipFuncs
			case "ShowFilenames":
				return !flags.SkipFilenames
			case "ShowConsts": 
				return !flags.SkipConsts
			case "ShowEmpty": 
				return flags.SkipEmpty
			case "ShowAll": 
				return !flags.SkipAll
			}
		panic("Invalid flag name " + flag_name)
	}
}