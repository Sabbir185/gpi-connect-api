package code

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

const numberBytes = "0123456789"

func GenerateRandomNumber(n int) (string, error) {
	b := make([]byte, n)
	max := big.NewInt(int64(len(numberBytes)))
	for i := range b {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = numberBytes[index.Int64()]
	}
	return string(b), nil
}

// HashOTP binds the code to its identity so a leaked digest can't be replayed for another number.
func HashOTP(secret, identity, otp string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(identity + ":" + otp))
	return hex.EncodeToString(mac.Sum(nil))
}

func EqualHash(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}
