package main

//TODO: Break up some components of the TUI into their own go files/directories
//TODO: Add legend for keybinds

import (
	"fmt"
	"os"

	"github.com/Zachkp/GoMail/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(models.CreateTable(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
