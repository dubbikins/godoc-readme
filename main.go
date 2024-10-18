package main

import "github.com/dubbikins/godoc-readme/cmd"

//go:generate go run main.go --skip-imports --skip-files


func main() {
	cmd.Execute()
}
