package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/12jihan/afk-x/internal/ui"
)

var version = "0.1.0-dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("afk-x %s\n", version)
		os.Exit(0)
	}

	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "afk-x: %v\n", err)
		os.Exit(1)
	}
}
