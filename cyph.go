package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding"
	"encoding/base32"
)

func encode_base32(data []byte) string {
	return base32.StdEncoding.EncodeToString(data);
}

func decode_base32(s string) ([]byte, error) {
    return base32.StdEncoding.DecodeString(s)
}

func generateNonce(data string) []byte {
	hash := sha512.Sum512([]byte(data))
	return hash[:12]
}

func hashgen(data string) string {
	if config_.hash_keys {
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

		return base32.StdEncoding.EncodeToString(hashkey.Sum(nil))
	} else {
		return base32.StdEncoding.EncodeToString([]byte(encrypt(data)));
	}
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
