package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Colors
var white string = "#FFFFFF"
var darkGray string = "#3C3C3C"
var green string = "#a6e3a1"

// Placeholder for dynamic table sizing
var placholderWidth int = 1

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(darkGray))

type model struct {
	table  table.Model
	width  int
	height int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	//TODO: Fix dynamic windo size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetColumns(createColumns(m.width))
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Email from: %s", m.table.SelectedRow()[0]),
				tea.Printf("Received on: %s", m.table.SelectedRow()[1]),
				tea.Printf("Width: %d", m.width),
				tea.Printf("Height: %d", m.height),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func createColumns(width int) []table.Column {
	// Reserve a little padding (e.g., for borders or spacing)
	totalWidth := width - 10

	// Allocate width proportionally
	senderW := 25
	dateW := 10
	timeW := 10
	messageW := totalWidth - senderW - dateW - timeW

	// Don't let message width go negative
	if messageW < 20 {
		messageW = 20
	}

	return []table.Column{
		{Title: "Sender", Width: senderW},
		{Title: "Date", Width: dateW},
		{Title: "Time", Width: timeW},
		{Title: "Message", Width: messageW},
	}
}

func main() {

	/*columns := []table.Column{
		{Title: "Sender", Width: 25},   // ex: johnsmith@email.com
		{Title: "Date", Width: 10},     // ex: 10/20/2025
		{Title: "Time", Width: 10},     // ex: 10:00AM
		{Title: "Message", Width: 125}, // Contents of email -- Show preview ??
	}*/

	columns := createColumns(placholderWidth)

	//TODO: Populate from email account, then use Bat to display whole email
	rows := []table.Row{
		{"Johndoe@email.com", "10/20/2025", "8:30AM", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
		{"Jasonk@email.com", "08/01/2025", "3:00PM", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(green)).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color(white)).
		Background(lipgloss.Color(darkGray)).
		Bold(true)

	t.SetStyles(s)

	m := model{
		table: t,
		width: placholderWidth,
	}

	//TODO: add back in tea.NewProgram(m,  tea.WithAltScreen())
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}

}
