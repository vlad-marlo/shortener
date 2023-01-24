package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
)

var (
	// encryptor ...
	encryptor *Encryptor

	// once
	once sync.Once
)

// Encryptor ...
type Encryptor struct {
	nonce []byte
	GCM   cipher.AEAD
}

// generateRandom byte slice
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("rand read: %w", err)
	}
	return b, nil
}

// EncodeUUID ...
func (e *Encryptor) EncodeUUID(uuid string) string {
	src := []byte(uuid)
	dst := e.GCM.Seal(nil, e.nonce, src, nil)
	return hex.EncodeToString(dst)
}

// DecodeUUID ...
func (e *Encryptor) DecodeUUID(uuid string, to *string) error {
	dst, err := hex.DecodeString(uuid)
	if err != nil {
		return fmt.Errorf("hex decode: %v", err)
	}

	src, err := e.GCM.Open(nil, e.nonce, dst, nil)
	if err != nil {
		return fmt.Errorf("gcm open: %v", err)
	}

	*to = string(src)
	return nil
}

// Get ...
func Get() *Encryptor {
	once.Do(func() {
		key, err := generateRandom(aes.BlockSize)
		if err != nil {
			log.Fatalf("generate key: %v", err)
		}

		aesBlock, err := aes.NewCipher(key)
		if err != nil {
			log.Fatalf("initialize cipher: %v", err)
		}

		aesGCM, err := cipher.NewGCM(aesBlock)
		if err != nil {
			log.Fatalf("initialize GCM encryptor: %v", err)
		}

		nonce, err := generateRandom(aesGCM.NonceSize())
		if err != nil {
			log.Fatalf("initialize GCM nonce: %v", err)
		}

		encryptor = &Encryptor{
			nonce: nonce,
			GCM:   aesGCM,
		}
	})
	return encryptor
}
