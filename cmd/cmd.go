package cmd

import (
	"fmt"
	"os"

	"github.com/dubbikins/godoc-readme/godoc_readme"
	"github.com/dubbikins/godoc-readme/godoc_readme/template_functions"
	"github.com/spf13/cobra"
)

var recursive bool = true
var template_filename string
var confirm_updates bool
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
		"confirm-updates","c", false,
		"Use this flag to confirm overwriting existing README.md files. The default behaviour is to overwrite the file without confirmation. Confirmation also gives you the option to view the diff between the existing and generated file before overwriting it.",
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
			if !recursive {
				ro.DirPattern = ro.Dir
			}
			if template_filename != "" {
				ro.TemplateFile = template_filename
			}
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