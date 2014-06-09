package eauth

import (
	"testing"
	"time"
)

func TestCreateLink(t *testing.T) {
	then := time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	now := then.Add(time.Minute)

	admin := User{Id: 1}

	link := createLink(localhost, admin, then)
	x := `http://localhost:8008/auth/1/1388534400/JPxBnXHa0zn3sYX7sjcbtlIPUQI=`
	if link != x {
		t.Errorf("Unexpected link from CreateLink(): %s", link)
	}
	if !isValid(localhost, admin, then, now, "JPxBnXHa0zn3sYX7sjcbtlIPUQI=") {
		t.Errorf("CreateLink() returned an invalid link")
	}

	// TODO Just with configurable link expiration
	// It is currently set to 1 hour
	now = now.Add(time.Hour)
	if isValid(localhost, admin, then, now, "JPxBnXHa0zn3sYX7sjcbtlIPUQI=") {
		t.Error("isValid() claimed a link was valid when it had expired")
	}
}
