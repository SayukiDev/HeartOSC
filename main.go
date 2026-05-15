package main

import (
	"HeartOSC/app"
	"HeartOSC/common/title"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	title.SetTitle("Heart OSC")
	p := tea.NewProgram(app.NewModel())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
