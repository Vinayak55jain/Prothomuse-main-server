package utils
import (
	"crypto/rand"
	"encoding/base64"
)
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil { // this will generate 32 bytes of random data in the byte slice.just polluting the memory.
		return "", err
	}
	return "pk_" + base64.URLEncoding.EncodeToString(bytes)[:40], nil // it create the 40 character long api key.
}