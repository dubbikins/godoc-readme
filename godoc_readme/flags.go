package godoc_readme

/* RenderFlags can be used to turn on and off rendering of different sections in the README.md file.

The bitmask values are as follows:

| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | ... | 32 |
|---|---|---|---|---|---|---|---|---|----|-----|----|
| RenderTypes | RenderFuncs | RenderMethods | RenderVars | RenderConsts | RenderExamples | RenderAlerts | TBD | TBD | TBD | TBD | RenderAll (default) |
*/
const (
	RenderTypes RenderFlag = 1 << iota
	RenderFuncs RenderFlag = 2 << iota
	RenderMethods
	RenderVars
	RenderConsts
	RenderExamples
	RenderAlerts
	RenderAll = ^RenderFlag(0)	
)

/* RenderFlags can be used to turn on and off rendering of different sections in the README.md file.

The bitmask values are as follows:

| 1 | 2 | 3 | 4 | 5 | 6 | 7 | ... | 32 |
|---|---|---|---|---|---|---|-----|----|
| Types | Funcs | TypeMethods | Vars | Consts | Examples | Alerts | TBD | RenderAll (default) |

For example, to render only the types and functions in the README.md file, you would set the 1st and 2nd bits, i.e `0000 0011` or `RenderTypes | RenderFuncs`

*/
type RenderFlag uint32 

// IsSet returns true if the flag is set in the RenderFlags
func (f RenderFlag) IsSet(flag RenderFlag) bool{
	return f & flag == f
}


