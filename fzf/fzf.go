// fzf/fzf.go
package fzf

import (
	"strings"

	"github.com/Zachkp/GoMail/email"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// FuzzySearchEmails searches through emails using fuzzy matching
func FuzzySearchEmails(emails []email.Email, query string) []email.Email {
	if query == "" {
		return emails
	}

	var results []email.Email
	var searchTargets []string

	// Create searchable strings from email data
	for _, e := range emails {
		searchTarget := strings.ToLower(e.From + " " + e.Subject + " " + e.Body)
		searchTargets = append(searchTargets, searchTarget)
	}

	// Find matches using fuzzy search
	matches := fuzzy.RankFindFold(strings.ToLower(query), searchTargets)

	// Sort by rank and return matching emails
	for _, match := range matches {
		if match.Distance <= len(query)*2 { // Adjust threshold as needed
			results = append(results, emails[match.OriginalIndex])
		}
	}

	return results
}

// FuzzySearchEmailsByField searches emails by specific field
func FuzzySearchEmailsByField(emails []email.Email, query string, field string) []email.Email {
	if query == "" {
		return emails
	}

	var results []email.Email
	var searchTargets []string

	// Create searchable strings based on field
	for _, e := range emails {
		var searchTarget string
		switch field {
		case "from":
			searchTarget = strings.ToLower(e.From)
		case "subject":
			searchTarget = strings.ToLower(e.Subject)
		case "body":
			searchTarget = strings.ToLower(e.Body)
		default:
			searchTarget = strings.ToLower(e.From + " " + e.Subject + " " + e.Body)
		}
		searchTargets = append(searchTargets, searchTarget)
	}

	// Find matches
	matches := fuzzy.RankFindFold(strings.ToLower(query), searchTargets)

	for _, match := range matches {
		if match.Distance <= len(query)*2 {
			results = append(results, emails[match.OriginalIndex])
		}
	}

	return results
}
