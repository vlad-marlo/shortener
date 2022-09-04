package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

// generateRandom byte slice
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

type Encryptor struct {
	nonce []byte
	GCM   cipher.AEAD
}

func NewEncryptor() (*Encryptor, error) {
	key, err := generateRandom(aes.BlockSize)
	if err != nil {
		return nil, err
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce, err := generateRandom(aesGCM.NonceSize())
	if err != nil {
		return nil, err
	}

	return &Encryptor{
		nonce: nonce,
		GCM:   aesGCM,
	}, nil
}

func (e *Encryptor) EncodeUUID(uuid string) string {
	src := []byte(uuid)
	dst := e.GCM.Seal(nil, e.nonce, src, nil)
	return hex.EncodeToString(dst)
}

func (e *Encryptor) DecodeUUID(uuid string) (string, error) {
	dst, err := hex.DecodeString(uuid)
	if err != nil {
		return "", err
	}
	src, err := e.GCM.Open(nil, e.nonce, dst, nil)
	if err != nil {
		return "", err
	}
	return string(src), nil
}
