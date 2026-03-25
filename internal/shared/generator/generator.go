package generator

import (
	"fmt"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bytesForCode = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int, charSet string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

func GeneratePassword() string {
	return randStringBytes(12, letterBytes)
}

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GeneratePickupCode() string {
	return randStringBytes(8, bytesForCode)
}
