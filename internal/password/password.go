package password

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	iterations = 600_000
	keyLength  = 32
)

func Hash(password string) (string, error) {
	salt := make([]byte, 16)
	rand.Read(salt)

	key, err := pbkdf2.Key(sha256.New, password, salt, iterations, keyLength)
	if err != nil {
		return "", fmt.Errorf("hachage du mot de passe : %w", err)
	}
	return hex.EncodeToString(salt) + "$" + hex.EncodeToString(key), nil
}

func Check(password, stored string) bool {
	encodedSalt, hashedPwd, found := strings.Cut(stored, "$")
	if !found {
		return false
	}

	salt, err := hex.DecodeString(encodedSalt)
	if err != nil {
		return false
	}

	storedPwd, err := hex.DecodeString(hashedPwd)
	if err != nil {
		return false
	}

	givenPwd, err := pbkdf2.Key(sha256.New, password, salt, iterations, keyLength)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(storedPwd, givenPwd) == 1
}
