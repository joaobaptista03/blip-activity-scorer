package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Repository Activity Scorer started")
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Stub pipeline to be filled in subsequent commits
	return nil
}
