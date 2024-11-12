/*
Godoc Readme CLI

Generate README.md file for your go project using comments you already write

Usage:
  godoc-readme [flags]

Flags:
  -c, --confirm          Use this flag to confirm overwriting existing README.md files. The default behaviour is to overwrite the file without confirmation. Confirmation also gives you the option to view the diff between the existing and generated file before overwriting it.
  -h, --help             help for godoc-readme
  -p, --package string   Specify the pattern for matching packages to generate the README.md files for. Default '' will match current package only
  -r, --recursive        If set, recursively search for go packages in the directory and generate a README.md for each package; Default will only create a Readme for the package found in the current directory
      --skip-all         Skips generating all sections besides the package documentation
      --skip-consts      Shows generating the consts section
      --skip-empty       Skips generating any type, func, var, const, or method that does not have a doc string
      --skip-examples    Skips generating the examples section
      --skip-filenames   Skips generating the files section
      --skip-funcs       Skips generating the functions section
      --skip-imports     Skips generating the imports section
      --skip-types       Skips generating the types section
      --skip-vars        Skips generating the vars section
*/
package cmd
