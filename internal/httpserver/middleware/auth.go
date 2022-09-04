package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Encryptor struct {
	nonce []byte
	GCM   cipher.AEAD
}

const (
	UserIDCookieName   = "user"
	UserCTXName        = "user_in_context"
	UserIDDefaultValue = "default_user"
)

var encryptor *Encryptor

// generateRandom byte slice
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// NewEncryptor ...
func NewEncryptor() error {
	if encryptor != nil {
		return nil
	}

	key, err := generateRandom(aes.BlockSize)
	if err != nil {
		return err
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return err
	}

	nonce, err := generateRandom(aesGCM.NonceSize())
	if err != nil {
		return err
	}

	encryptor = &Encryptor{
		nonce: nonce,
		GCM:   aesGCM,
	}

	return nil
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
		return err
	}

	src, err := e.GCM.Open(nil, e.nonce, dst, nil)
	if err != nil {
		return err
	}

	*to = string(src)
	return nil
}

// AuthMiddleware ...
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rawUserID string

		if err := NewEncryptor(); err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCTXName, UserIDDefaultValue)))
			return
		}

		if user, err := r.Cookie(UserIDCookieName); err != nil {
			rawUserID = uuid.New().String()
		} else if err = encryptor.DecodeUUID(user.Value, &rawUserID); err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCTXName, UserIDDefaultValue)))
			return
		}

		encoded := encryptor.EncodeUUID(rawUserID)
		c := &http.Cookie{
			Name:  UserIDCookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, c)
		log.Print(rawUserID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCTXName, rawUserID)))
	})
}
