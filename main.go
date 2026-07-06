// Package main is the entrypoint for the uts CLI.
package main

import (
	"fmt"
	"os"

	"github.com/y3owk1n/uts/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
