package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// Prints a random secret for JWT_SECRET or OTP_SECRET.
func main() {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("generate secret: %v", err)
	}
	fmt.Println(base64.RawURLEncoding.EncodeToString(b))
}
