package utils

import (
	"math/rand"
	"time"
)

// Reference by "https://www.calhoun.io/creating-random-strings-in-go/"
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func getRandStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetRandString : 랜덤 문자열 획득
func GetRandString(length int) string {
	return getRandStringWithCharset(length, charset)
}
