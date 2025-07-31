package styles

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	White    string = "#FFFFFF"
	DarkGray string = "#3C3C3C"
	Green    string = "#a6e3a1"
)

// Layout
var (
	PlaceholderWidth = 1 // placeholder for dynamic table sizing
)

// TODO: Maybe expand??
// Base Styles
var (
	BaseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(DarkGray))
)
