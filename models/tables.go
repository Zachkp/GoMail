package models

import (
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func CreateColumns(width int) []table.Column {
	// Reserve a little padding (e.g., for borders or spacing)
	totalWidth := width - 10

	// Allocate width proportionally
	senderWidth := 25
	dateWidth := 10
	timeWidth := 10
	messageWidth := totalWidth - senderWidth - dateWidth - timeWidth

	// Don't let message width go negative
	if messageWidth < 20 {
		messageWidth = 20
	}

	return []table.Column{
		{Title: "Sender", Width: senderWidth},
		{Title: "Date", Width: dateWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Message", Width: messageWidth},
	}
}

func CreateTable() model {
	columns := CreateColumns(styles.PlaceholderWidth)

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

	m := model{
		table: t,
		width: styles.PlaceholderWidth,
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(styles.Green)).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color(styles.White)).
		Background(lipgloss.Color(styles.DarkGray)).
		Bold(true)

	t.SetStyles(s)

	return m
}
