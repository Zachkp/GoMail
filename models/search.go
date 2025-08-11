// models/search.go
package models

import (
	"strings"

	"github.com/Zachkp/GoMail/email"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Add these fields to your existing model struct in models.go
type SearchState struct {
	isSearching    bool
	searchInput    textinput.Model
	originalEmails []email.Email
	filteredEmails []email.Email
}

// Initialize search functionality
func InitSearch() SearchState {
	ti := textinput.New()
	ti.Placeholder = "Search emails..."
	ti.CharLimit = 100
	ti.Width = 50

	return SearchState{
		isSearching:    false,
		searchInput:    ti,
		originalEmails: []email.Email{},
		filteredEmails: []email.Email{},
	}
}

// Simple fuzzy search function (in case lithammer doesn't work)
func simpleFuzzyMatch(query, target string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	target = strings.ToLower(target)

	if query == "" {
		return true
	}

	// Simple substring matching for now
	return strings.Contains(target, query)
}

// Update search input and filter results
func (s *SearchState) UpdateSearch(value string, allEmails []email.Email) {
	if value == "" {
		s.filteredEmails = s.originalEmails
		return
	}

	var results []email.Email
	query := strings.ToLower(strings.TrimSpace(value))

	for _, email := range s.originalEmails {
		searchTarget := strings.ToLower(email.From + " " + email.Subject + " " + email.Body)

		// Try fuzzy search first, fallback to simple matching
		var matches bool

		// Use fuzzy search if available
		searchTargets := []string{searchTarget}
		fuzzyMatches := fuzzy.RankFindFold(query, searchTargets)
		if len(fuzzyMatches) > 0 && fuzzyMatches[0].Distance <= len(query)*2 {
			matches = true
		} else {
			// Fallback to simple substring matching
			matches = simpleFuzzyMatch(query, searchTarget)
		}

		if matches {
			results = append(results, email)
		}
	}

	s.filteredEmails = results
}

// Toggle search mode
func (s *SearchState) ToggleSearch(emails []email.Email) {
	if !s.isSearching {
		// Entering search mode
		s.originalEmails = make([]email.Email, len(emails))
		copy(s.originalEmails, emails)
		s.filteredEmails = make([]email.Email, len(emails))
		copy(s.filteredEmails, emails)
		s.searchInput.Focus()
		s.isSearching = true
	} else {
		// Exiting search mode
		s.searchInput.Blur()
		s.searchInput.SetValue("")
		s.filteredEmails = s.originalEmails
		s.isSearching = false
	}
}

// Render search bar
func (s *SearchState) RenderSearchBar() string {
	if !s.isSearching {
		return ""
	}

	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	return searchStyle.Render("🔍 " + s.searchInput.View())
}
