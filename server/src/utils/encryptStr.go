package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func EncryptStr(str string) string {
	bytes := []byte(str)

	// // Using bcrypt with random salt
	// encrypted, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)

	hasher := sha256.New()
	hasher.Write(bytes)

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha
}
