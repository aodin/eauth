package eauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

// Links can be reused within the expiration, but I'm okay with that
func CreateLink(c Config, uid int64) string {
	return createLink(c, uid, time.Now())
}

// TODO How long is a link valid for?
func createLink(c Config, uid int64, now time.Time) string {
	d := createDigest(c, uid, now)

	// TODO Use url.URL?
	protocol := "http"
	if c.Https {
		protocol += "s"
	}

	// TODO Allow the auth url to be set

	// Create a link from the uid, timestamp, and hmac digest
	return fmt.Sprintf(
		"%s://%s/auth/%s/%s/%s",
		protocol,
		c.Domain,
		fmt.Sprint(uid),
		fmt.Sprint(now.Unix()),
		base64.URLEncoding.EncodeToString(d),
	)
}

func createDigest(c Config, uid int64, now time.Time) []byte {
	// TODO Encode the user id in another base?
	// TODO This byte casts are silly
	h := hmac.New(sha1.New, []byte(c.Secret))
	h.Write([]byte(fmt.Sprint(uid)))
	h.Write([]byte(fmt.Sprint(now.Unix())))
	return h.Sum(nil)
}

func isValid(c Config, uid int64, then, now time.Time, given string) bool {
	// Test that the given digest matches the expected
	d := createDigest(c, uid, then)
	expected := base64.URLEncoding.EncodeToString(d)

	// Just a proxy for subtle.ConstantTimeCompare
	if !hmac.Equal([]byte(expected), []byte(given)) {
		return false
	}

	// TODO Link expiration should be configurable
	validFor := time.Hour

	// TODO Confirm that the user is valid?
	return now.Before(then.Add(validFor))
}
