package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
)

const (
	stateLength        = 96
	codeVerifierLength = 32
)

type Session struct {
	state         string
	codeVerifier  string
	codeChallenge string
}

func newSession() *Session {
	bytes := getRandomStringBytes(stateLength)
	state := base64.RawURLEncoding.EncodeToString(bytes)

	bytes = getRandomStringBytes(codeVerifierLength)
	codeVerifier := base64.RawURLEncoding.EncodeToString(bytes)

	hashBytes := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hashBytes[:])

	return &Session{
		state:         state,
		codeVerifier:  codeVerifier,
		codeChallenge: codeChallenge,
	}
}

func getRandomStringBytes(n int) []byte {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	runes := make([]rune, n)
	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}

	return []byte(string(runes))
}
