package pcke

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"
)

// They said it's best to use crypto/rand.
// The specification explicitly states these are the only allowed characters
func generateRandomString(n int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_.-~"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return nil, err
		}
		ret[i] = letters[num.Int64()]
	}

	return ret, nil
}

func sha256UrlEncode(b []byte) string {
	res := sha256.Sum256(b)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(res[:]), "=")
}

func generateRandomState() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
