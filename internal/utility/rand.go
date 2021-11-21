package utility

import (
	"crypto/rand"
	"math/big"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandomString(length int) *string {
	bts := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(62)))
		if err != nil {
			return nil // Let the program crash
		}

		bts[i] = letters[num.Int64()]
	}

	str := string(bts)

	return &str
}
