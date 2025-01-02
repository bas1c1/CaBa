package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

func generateNonce(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:12]
}

func encrypt(data string) string {
	block, err := aes.NewCipher(config_.passkey)
	if err != nil {
		throw(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		throw(err)
	}

	nonce := generateNonce(data)

	return string(gcm.Seal(nonce, nonce, []byte(data), nil))
}

func decrypt(line string) string {
	block, err := aes.NewCipher(config_.passkey)
	if err != nil {
		throw(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		throw(err)
	}

	nonceSize := gcm.NonceSize()
	if len(line) < nonceSize {
		throw("line too short")
	}

	nonce, ciphertext := []byte(line[:nonceSize]), []byte(line[nonceSize:])

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		throw(err)
	}

	return string(plaintext)
}
