package email

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/joho/godotenv"
)

// Email represents a simplified email record with body text
type Email struct {
	From    string
	Subject string
	Date    string
	Body    string // plain text body
}

// FetchLatestEmails connects to the IMAP server and returns the newest N emails (newest-first) with plain text bodies
func FetchLatestEmails(limit uint32) ([]Email, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_IMAP_HOST")
	port := os.Getenv("EMAIL_IMAP_PORT")

	if username == "" || password == "" || host == "" || port == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	// Connect to server
	c, err := client.DialTLS(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(username, password); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}

	if mbox.Messages == 0 {
		return []Email{}, nil
	}

	// Determine range for newest N messages
	var from uint32
	if mbox.Messages > limit {
		from = mbox.Messages - limit + 1
	} else {
		from = 1
	}
	to := mbox.Messages

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// Fetch both envelope and body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}

	messages := make(chan *imap.Message, limit)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var emails []Email
	for msg := range messages {
		fromAddr := ""
		if len(msg.Envelope.From) > 0 {
			fromAddr = msg.Envelope.From[0].Address()
		}

		// Extract the body
		body := ""
		if r := msg.GetBody(section); r != nil {
			mr, err := mail.CreateReader(r)
			if err == nil {
				// Loop through all parts
				for {
					p, err := mr.NextPart()
					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}

					switch h := p.Header.(type) {
					case *mail.InlineHeader:
						ct := h.Get("Content-Type")
						if strings.HasPrefix(ct, "text/plain") {
							b, _ := io.ReadAll(p.Body)
							body = string(b)
						}
					}
				}
			}
		}

		emails = append(emails, Email{
			From:    fromAddr,
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date.Format("2006-01-02 15:04:05"),
			Body:    body,
		})
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// Reverse slice for newest-first
	for i, j := 0, len(emails)-1; i < j; i, j = i+1, j-1 {
		emails[i], emails[j] = emails[j], emails[i]
	}

	return emails, nil
}
