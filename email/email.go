package email

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

// Email represents a simplified email record with body text (plain or HTML as-is)
type Email struct {
	From    string
	Subject string
	Date    string
	Body    string
}

func htmlToPlainText(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}

	var buf strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "style" || n.Data == "script" || n.Data == "head" || n.Data == "meta" || n.Data == "link") {
			return
		}

		if n.Type == html.TextNode {
			// Clean up whitespace (collapse multiple spaces)
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				buf.WriteString(text)
				buf.WriteString(" ")
			}
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "div", "br", "tr":
				buf.WriteString("\n")
			case "li":
				buf.WriteString("\n- ")
			case "a":
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return strings.TrimSpace(buf.String())
}

func collapseBlankLines(text string) string {
	re := regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n")
	return text
}

func FetchLatestEmails(limit uint32) ([]Email, error) {
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

	// Connect and login
	c, err := client.DialTLS(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer c.Logout()

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

	var from uint32
	if mbox.Messages > limit {
		from = mbox.Messages - limit + 1
	} else {
		from = 1
	}
	to := mbox.Messages

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

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

		body := ""

		if r := msg.GetBody(section); r != nil {
			mr, err := mail.CreateReader(r)
			if err == nil {
				var htmlBody, plainBody string

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
						b, _ := io.ReadAll(p.Body)
						content := string(b)
						if strings.HasPrefix(ct, "text/html") && htmlBody == "" {
							htmlBody = content
						} else if strings.HasPrefix(ct, "text/plain") && plainBody == "" {
							plainBody = content
						}
					}
				}

				if htmlBody != "" {
					plainText := htmlToPlainText(htmlBody)
					cleanText := collapseBlankLines(plainText)
					body = cleanText

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

	for i, j := 0, len(emails)-1; i < j; i, j = i+1, j-1 {
		emails[i], emails[j] = emails[j], emails[i]
	}

	return emails, nil
}
