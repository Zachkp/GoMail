package models

import (
	"fmt"

	"github.com/Zachkp/GoMail/email"
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table         table.Model
	width, height int
	emails        []email.Email
	viewingEmail  bool
	selectedEmail email.Email
	emailViewport viewport.Model
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

		if m.viewingEmail {
			containerHeight := m.height - 6
			headerHeight := 4
			viewportHeight := containerHeight - headerHeight - 2
			if viewportHeight < 5 {
				viewportHeight = 5
			}
			m.emailViewport.Width = m.width - 8
			m.emailViewport.Height = viewportHeight
		}

		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, CommonKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, CommonKeys.Select):
			if !m.viewingEmail {
				selectedRow := m.table.Cursor()
				if selectedRow >= 0 && selectedRow < len(m.emails) {
					m.selectedEmail = m.emails[selectedRow]
					m.viewingEmail = true

					containerHeight := m.height - 6
					headerHeight := 4
					viewportHeight := containerHeight - headerHeight - 2
					if viewportHeight < 5 {
						viewportHeight = 5
					}

					m.emailViewport = viewport.New(m.width-8, viewportHeight)
					m.emailViewport.SetContent(m.selectedEmail.Body)
				}
			}

		case key.Matches(msg, CommonKeys.Back):
			if m.viewingEmail {
				m.viewingEmail = false
			}

		case key.Matches(msg, CommonKeys.Up):
			if m.viewingEmail {
				m.emailViewport.LineUp(1)
				return m, nil
			} else {
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}

		case key.Matches(msg, CommonKeys.Down):
			if m.viewingEmail {
				m.emailViewport.LineDown(1)
				return m, nil
			} else {
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
		containerHeight := m.height - 6

		headerContent := fmt.Sprintf(
			"From: %s\nDate: %s\nTime: %s\nSubject: %s\n",
			m.selectedEmail.From,
			func() string {
				if len(m.selectedEmail.Date) >= 10 {
					return m.selectedEmail.Date[:10]
				}
				return m.selectedEmail.Date
			}(),
			func() string {
				if len(m.selectedEmail.Date) >= 16 {
					return m.selectedEmail.Date[11:16]
				}
				return ""
			}(),
			m.selectedEmail.Subject,
		)

		headerView := lipgloss.NewStyle().
			Bold(true).
			Padding(0, 0, 1, 0).
			Render(headerContent)

		emailBodyView := m.emailViewport.View()

		emailView := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(styles.Green)).
			Padding(1, 2).
			Width(m.width - 8).
			Height(containerHeight).
			Render(lipgloss.JoinVertical(lipgloss.Left, headerView, emailBodyView))

		helpView := CommonHelp.View(CommonKeys)

		return lipgloss.JoinVertical(lipgloss.Center, emailView, helpView)
	}

	// Default table view:
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
