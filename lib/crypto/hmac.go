package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

// HMac crypto payload
func HMac(body interface{}, secret string) string {
	message, _ := json.Marshal(body)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(message)
	// to lowercase hexits
	hex.EncodeToString(hash.Sum(nil))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC interface{}, key string) bool {
	if message == nil || messageMAC == nil || key == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message.(string)))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal([]byte(messageMAC.(string)), expectedMAC)
}
