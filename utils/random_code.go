package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var codes []byte = make([]byte, length)

	for i := 0; i < length; i++ {
		codes[i] = uint8(48 + r.Intn(10))
	}

	return string(codes)
}
