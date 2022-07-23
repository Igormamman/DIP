package tools

import (
	"crypto/sha256"
	"fmt"
)

func GenerateSHA256(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func StringGenerateSHA256(data string) string {
	return fmt.Sprintf("%x",sha256.Sum256([]byte(data)))
}


