package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding"
	"encoding/hex"
)

func generateNonce(data string) []byte {
	hash := sha512.Sum512([]byte(data))
	return hash[:12]
}

func hashgen(data string) string {
	hashkey := sha512.New()
	hashkey.Write([]byte(data))

	marshaler, ok := hashkey.(encoding.BinaryMarshaler)
	if !ok {
		caba_err("first does not implement encoding.BinaryMarshaler")
	}
	_, err := marshaler.MarshalBinary()
	if err != nil {
		caba_err("unable to marshal hash:")
	}

	return hex.EncodeToString(hashkey.Sum(nil))
}

func encrypt(data string) string {
	block, err := aes.NewCipher(config_.passkey)
	if err != nil {
		throw(err)
		return ""
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		throw(err)
		return ""
	}

	nonce := generateNonce(data)

	return string(gcm.Seal(nonce, nonce, []byte(data), nil))
}

func decrypt(line []byte) string {
	block, err := aes.NewCipher(config_.passkey)
	if err != nil {
		throw(err)
		return ""
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		throw(err)
		return ""
	}

	nonceSize := gcm.NonceSize()
	if len(line) < nonceSize {
		throw("line too short")
		return ""
	}

	nonce, ciphertext := []byte(line[:nonceSize]), []byte(line[nonceSize:])

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		throw(err)
		return ""
	}

	return string(plaintext)
}
