// Package main is the entrypoint for the genman CLI.
//
//nolint:mnd
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra/doc"
	"github.com/y3owk1n/uts/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: genman <output-dir>\n")
		os.Exit(1)
	}

	outputDir := os.Args[1]

	err := os.MkdirAll(outputDir, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	if epoch := os.Getenv("SOURCE_DATE_EPOCH"); epoch != "" {
		sec, err := strconv.ParseInt(epoch, 10, 64)
		if err == nil {
			now = time.Unix(sec, 0).UTC()
		}
	}

	header := &doc.GenManHeader{
		Title:   "UTS",
		Section: "1",
		Date:    &now,
		Manual:  "uts Manual",
		Source:  "uts " + cmd.Version,
	}

	err = doc.GenManTree(cmd.RootCmd, header, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating man pages: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Man pages generated in %s/\n", outputDir) //nolint:forbidigo
}
