package tokens

import (
	"math/rand"
	"time"
)

const LENGTH = 16

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateString(length int) string {
	return stringWithCharset(length, charset)
}

func GenerateToken() string {
	return generateString(LENGTH)
}
