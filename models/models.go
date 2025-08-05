package models

import (
	"fmt"

	"github.com/Zachkp/GoMail/email"
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table         table.Model
	width, height int
	emails        []email.Email // store all email data
	viewingEmail  bool          // are we viewing a single email?
	selectedEmail email.Email   // currently selected email
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetColumns(CreateColumns(m.width - 20))
		m.table.SetHeight(m.height - 20)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, CommonKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, CommonKeys.Select):
			if !m.viewingEmail {
				// Enter email view
				selectedRow := m.table.Cursor()
				if selectedRow >= 0 && selectedRow < len(m.emails) {
					m.selectedEmail = m.emails[selectedRow]
					m.viewingEmail = true
				}
			} else {
				// Already in email view? Ignore select
			}

		case key.Matches(msg, CommonKeys.Back):
			if m.viewingEmail {
				// Go back to table view
				m.viewingEmail = false
			}

		case key.Matches(msg, CommonKeys.Up):
			if !m.viewingEmail {
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}

		case key.Matches(msg, CommonKeys.Down):
			if !m.viewingEmail {
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		}
	}

	if !m.viewingEmail {
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.viewingEmail {
		// Render full email body view — keep padding and styling consistent

		emailStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(styles.Green)).
			Padding(1, 2).
			Width(m.width - 8).
			Height(m.height - 6) // adjust as needed

		emailContent := fmt.Sprintf(
			"From: %s\nDate: %s\nSubject: %s\n\n%s",
			m.selectedEmail.From,
			m.selectedEmail.Date,
			m.selectedEmail.Subject,
			m.selectedEmail.Body,
		)

		emailView := emailStyle.Render(emailContent)

		helpView := CommonHelp.View(CommonKeys)

		// Combine email body view + help just like your table layout
		return lipgloss.JoinVertical(lipgloss.Center, emailView, helpView)
	}

	// Default table view, exactly your current config:
	tableView := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(styles.Green)).
		Padding(0, 1)

	helpView := CommonHelp.View(CommonKeys)

	bordered := tableView.Render(m.table.View())

	padded := lipgloss.NewStyle().
		Padding(2, 4).
		Render(bordered)

	return lipgloss.JoinVertical(lipgloss.Center, padded, helpView)
}
