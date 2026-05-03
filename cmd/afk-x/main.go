package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/12jihan/afk-x/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
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
		log.Printf("afk-x: %v", err)
		os.Exit(1)
	}
}
