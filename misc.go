package libwekan

import (
	"math/rand"
	"time"
)

func newId() string {
	chars := "123456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"
	l := len(chars)
	var digits []byte
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 17; i++ {
		digit := rand.Intn(l)
		digits = append(digits, chars[digit])
	}
	return string(digits)
}
