package models

import (
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table  table.Model
	width  int
	height int
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

	//TODO: Add Logic for some of these
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, CommonKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, CommonKeys.Up):

		case key.Matches(msg, CommonKeys.Down):

		case key.Matches(msg, CommonKeys.Search):

		case key.Matches(msg, CommonKeys.Select):

		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
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
