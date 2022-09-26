package libwekan

import (
	"math/rand"
	"time"
)

func newId() string {
	return newIdN(17)
}

func newId6() string {
	return newIdN(6)
}

func newIdN(n int) string {
	chars := "123456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"
	l := len(chars)
	var digits []byte
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < n; i++ {
		digit := rand.Intn(l)
		digits = append(digits, chars[digit])
	}
	return string(digits)
}
