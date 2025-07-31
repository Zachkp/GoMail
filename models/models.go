package models

import (
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
		m.table.SetColumns(CreateColumns(m.width))
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
	return styles.BaseStyle.Render(m.table.View()) + "\n"
}
