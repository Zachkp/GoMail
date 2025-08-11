// Updated models/models.go - Add search integration
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

	// Add search functionality
	search SearchState
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

		// Update search input width
		m.search.searchInput.Width = m.width - 20

		return m, nil

	case tea.KeyMsg:
		// Handle search input first if we're searching
		if m.search.isSearching && !m.viewingEmail {
			switch {
			case key.Matches(msg, CommonKeys.Search): // Toggle search off
				m.search.ToggleSearch(m.emails)
				m.updateTableRows()
				return m, nil
			case key.Matches(msg, CommonKeys.Quit):
				return m, tea.Quit
			case msg.Type == tea.KeyEscape:
				m.search.ToggleSearch(m.emails)
				m.updateTableRows()
				return m, nil
			case msg.Type == tea.KeyEnter:
				// Exit search mode but keep results
				m.search.isSearching = false
				m.search.searchInput.Blur()
				return m, nil
			default:
				// Update search input
				m.search.searchInput, cmd = m.search.searchInput.Update(msg)
				m.search.UpdateSearch(m.search.searchInput.Value(), m.emails)
				m.updateTableRows()
				return m, cmd
			}
		}

		// Regular key handling
		switch {
		case key.Matches(msg, CommonKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, CommonKeys.Search):
			if !m.viewingEmail {
				m.search.ToggleSearch(m.emails)
				m.updateTableRows()
				return m, nil
			}

		case key.Matches(msg, CommonKeys.Select):
			if !m.viewingEmail {
				selectedRow := m.table.Cursor()
				currentEmails := m.getCurrentEmails()
				if selectedRow >= 0 && selectedRow < len(currentEmails) {
					m.selectedEmail = currentEmails[selectedRow]
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

	if !m.viewingEmail && !m.search.isSearching {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

// Helper function to get current emails (filtered or all)
func (m model) getCurrentEmails() []email.Email {
	if m.search.isSearching {
		return m.search.filteredEmails
	}
	return m.emails
}

// Helper function to update table rows
func (m *model) updateTableRows() {
	currentEmails := m.getCurrentEmails()
	var rows []table.Row

	for _, e := range currentEmails {
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

	m.table.SetRows(rows)
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

	// Default table view with search:
	tableView := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(styles.Green)).
		Padding(0, 1)

	helpView := CommonHelp.View(CommonKeys)

	// Build the view components
	var viewComponents []string

	// Add search bar if searching
	if searchBar := m.search.RenderSearchBar(); searchBar != "" {
		viewComponents = append(viewComponents, searchBar)
	}

	// Add table
	bordered := tableView.Render(m.table.View())
	padded := lipgloss.NewStyle().
		Padding(2, 4).
		Render(bordered)

	viewComponents = append(viewComponents, padded)

	// Add help
	viewComponents = append(viewComponents, helpView)

	return lipgloss.JoinVertical(lipgloss.Center, viewComponents...)
}
