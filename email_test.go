package eauth

import (
	"testing"
)

var testConfig = SMTPConfig{
	From:     "fred@example.com",
	Alias:    "Fred Flintstone",
	Host:     "example.com",
	Port:     587,
	User:     "Fred",
	Password: "yabbadabbado",
}

var x = `From: fred@example.com
To: <barney@example.com>, <wilma@example.com>
Subject: yabbadabbado
Content-Type: text/html; charset=UTF-8


yabbadabbado
`

func TestEmail(t *testing.T) {
	e := Email{
		From:    "fred@example.com",
		To:      Addresses{"barney@example.com", "wilma@example.com"},
		Subject: "yabbadabbado",
		Header:  map[string]string{"Content-Type": "text/html; charset=UTF-8"},
		Body:    "yabbadabbado",
	}
	out := e.String()
	if out != x {
		t.Errorf("Unexpected email string output: %s", out)
	}
}
