package main

import (
	"fmt"
	"os"
)

func main() {
	if err := BuildCLI().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
