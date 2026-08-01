package main

import (
	"fmt"
	"os"

	"github.com/sudabon/medit/internal/loader"
	"github.com/sudabon/medit/internal/ui"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: medit <file.md>")
	}

	doc, err := loader.Load(args[0])
	if err != nil {
		return err
	}

	return ui.Run(doc)
}
