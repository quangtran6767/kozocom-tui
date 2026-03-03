package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/config"
)

func main() {
	config.InitLogger()

	p := tea.NewProgram(newAppModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
