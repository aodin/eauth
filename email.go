package eauth

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

// Do not use port 465, it expects TLS from the start
// http://stackoverflow.com/a/11664176

type Sender interface {
	Send(c SMTPConfig, to []string, subject, body string) error
}

var msg = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}
{{ range $key, $value := .Header }}{{ $key }}: {{ $value }}
{{ end }}

{{ .Body }}
`
var emailTemplate = template.Must(template.New("email").Parse(msg))

type Addresses []string

func (a Addresses) String() string {
	escaped := make([]string, len(a))
	for i, address := range a {
		escaped[i] = fmt.Sprintf("<%s>", address)
	}
	return strings.Join(escaped, ", ")
}

// TODO Use the Values type for Header?
type Email struct {
	From    string
	To      Addresses
	Subject string
	Header  map[string]string
	Body    string
}

// TODO Prevent email lines from being over 78 characters?
func (email Email) String() string {
	b := new(bytes.Buffer)
	emailTemplate.Execute(b, email)
	// TODO What to do with errors?
	return b.String()
}

// DefaultSender implements the Email Sender interface
type DefaultSender struct{}

// Send will send an email on the DefaultSender
func (ds *DefaultSender) Send(c SMTPConfig, to Addresses, subject, body string) error {
	// Create the auth credentials using the given config
	auth := smtp.PlainAuth(
		"",
		c.User,
		c.Password,
		c.Host,
	)

	// Create the email
	// HTML and UTF-8 please
	mail := Email{
		From:    c.FromAddress(),
		To:      to,
		Subject: subject,
		Header:  map[string]string{"Content-Type": "text/html; charset=UTF-8"},
		Body:    body,
	}

	return smtp.SendMail(
		c.HostWithPort(),
		auth,
		c.From, // Do not use the alias!
		to,
		[]byte(mail.String()),
	)
}

// Default is an initialized instance of the DefaultSender
// No need to place this in init, there's no initialization behavior needed
var Default = &DefaultSender{}

// Send will send an email using the default Sender implementation.
func Send(c SMTPConfig, to []string, subject, body string) error {
	return Default.Send(c, to, subject, body)
}
