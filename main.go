package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/whichxjy/chords/ui"
)

func main() {
	if _, err := tea.NewProgram(&ui.Model{}).Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
