package models

import (
	"log"

	"github.com/Zachkp/GoMail/email"
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// CreateColumns sets up the column layout for the table
func CreateColumns(width int) []table.Column {
	// Allocate width proportionally
	senderWidth := 30
	dateWidth := 10
	timeWidth := 10
	messageWidth := width - senderWidth - dateWidth - timeWidth

	return []table.Column{
		{Title: "Sender", Width: senderWidth},
		{Title: "Date", Width: dateWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Message", Width: messageWidth},
	}
}

func CreateTable() model {
	columns := CreateColumns(styles.PlaceholderWidth)

	emails, err := email.FetchLatestEmails(25)
	if err != nil {
		log.Printf("Error fetching emails: %v", err)
		emails = []email.Email{}
	}

	var rows []table.Row
	for _, e := range emails {
		datePart := ""
		timePart := ""
		if len(e.Date) >= 10 {
			datePart = e.Date[:10]
		}
		if len(e.Date) >= 16 {
			timePart = e.Date[11:16]
		}

		rows = append(rows, table.Row{
			e.From,
			datePart,
			timePart,
			e.Subject,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	m := model{
		table:  t,
		width:  styles.PlaceholderWidth,
		emails: emails, // keep the full email data here
	}

	// styling...
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.ThickBorder()).
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
