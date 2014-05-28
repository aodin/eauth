package eauth

import (
	"crypto/hmac"
	"crypto/sha1"
)

// The validity of the reset link is determined by
// * The HMAC using the server's secret key
// * The unix timestamp in the message body
// * That the associated session does not have an assigned IP
func SaltedHMAC(salt, secret, data []byte) []byte {
	// Calculate the SHA1 digest of the secret + salt
	h := sha1.New()
	h.Write(salt)
	h.Write(secret)
	key := h.Sum(nil)

	// Create the HMAC
	hmacd := hmac.New(sha1.New, key)
	hmacd.Write(data)
	return hmacd.Sum(nil)
}

// func NewResetLink() string {
// 	salt := []byte("eauth.NewResetLink")
// }
