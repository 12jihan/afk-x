package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "0.1.0-dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("afk-x %s\n", version)
		os.Exit(0)
	}

	// Placeholder — Story 1.5 will replace this with tea.NewProgram(...)
	fmt.Println("afk-x: TUI not yet initialized")
}
