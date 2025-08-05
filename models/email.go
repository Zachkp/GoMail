package models

import (
	"time"
)

type ImapSmtpStore struct {
	ImapHost string
	ImapPort string
	SmtpHost string
	SmtpPort string
	Username string
	Password string
}

// Constructor function
func NewImapSmtpStore(username, password, imapHost, imapPort, smtpHost, smtpPort string) *ImapSmtpStore {
	return &ImapSmtpStore{
		Username: username,
		Password: password,
		ImapHost: imapHost,
		ImapPort: imapPort,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
	}
}

type Email struct {
	ID        string
	From      string
	To        string
	Subject   string
	Body      string
	Timestamp time.Time
	Read      bool
}

/*
type FileStore struct {
	path   string
	emails []Email
	mu     sync.Mutex
}
*/
