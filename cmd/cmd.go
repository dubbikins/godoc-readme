package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dubbikins/godoc-readme/godoc_readme"
	"github.com/dubbikins/godoc-readme/godoc_readme/template_functions"
	"github.com/spf13/cobra"
)

var recursive bool = true
var env string
var confirm_updates bool
var package_root string
// NOTE(flags): These Flags are used to determine which sections of the README.md file to generate

var flags template_functions.Flags = template_functions.Flags{
}

// Initializes the CLI flags/Arguments
func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&recursive, 
		"recursive", "r", false, 
		"If set, recursively search for go packages in the directory and generate a README.md for each package; Default will only create a Readme for the package found in the current directory",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&confirm_updates, 
		"confirm","c", false,
		"Use this flag to confirm overwriting existing README.md files. The default behaviour is to overwrite the file without confirmation. Confirmation also gives you the option to view the diff between the existing and generated file before overwriting it.",
	)
	rootCmd.PersistentFlags().StringVarP(
		&package_root, 
		"package","p", "",
		"Specify the pattern for matching packages to generate the README.md files for. Default '' will match current package only",
	)
	rootCmd.PersistentFlags().StringVarP(
		&env, 
		"env","e", "",
		"Specify the environment variables that should be passed to the build system. Example: 'GOOS=linux GOARCH=amd64'",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipExamples, 
		"skip-examples", false,
		"Skips generating the examples section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipFuncs, 
		"skip-funcs", false,
		"Skips generating the functions section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipConsts, 
		"skip-consts", false,
		"Shows generating the consts section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipImports, 
		"skip-imports", false,
		"Skips generating the imports section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipTypes, 
		"skip-types", false,
		"Skips generating the types section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipVars, 
		"skip-vars", false,
		"Skips generating the vars section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipFilenames, 
		"skip-filenames", false,
		"Skips generating the files section",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipEmpty, 
		"skip-empty", false,
		"Skips generating any type, func, var, const, or method that does not have a doc string",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flags.SkipAll, 
		"skip-all", false,
		"Skips generating all sections besides the package documentation",
	)
	
	// rootCmd.PersistentFlags().StringVarP(
	// 	&template_filename, 
	// 	"template", "t", "", 
	// 	"The template file to use for generating the README.md file",
	// )
	//rootCmd.Flags().BoolP("recursive", "r", true, "Recursively search for go packages in the directory and generate a README.md for each package")
}

// The root command for the CLI which passes the flags to the [godoc_readme package](../godoc_readme/README.md)
var rootCmd = &cobra.Command{
	Use:   "godoc-readme",
	Short: "Generate README.md file for your go project using comments you already write",
	Long:  `Generate README.md file for your go project using comments you already write`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(flags)
		if readme, err := godoc_readme.NewReadme(func(ro *godoc_readme.ReadmeOptions) {
				ro.PackageDir = package_root
				if recursive {
					ro.PackageDir = "./..."
				}
				ro.Env = strings.Split(env, "")
				ro.ConfirmUpdates = confirm_updates
				ro.Flags = flags
				
		}); err != nil {
				fmt.Println("err")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
		} else {
			if err = readme.Generate(); err != nil {
				fmt.Println("Generate err")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	},
}

// Execute runs the root command using the os.Args by default
// Optionally, you can pass in a list of arguments to run the command with
func Execute(args ...string) error{
	// if len(args) > 0 {
		
	// }
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		//slog.Error(err.Error())
		return err
	}
	return nil
}