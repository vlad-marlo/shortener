package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// types ...
type (
	// Encryptor ...
	Encryptor struct {
		nonce []byte
		GCM   cipher.AEAD
	}
	// cookieUserIDValueType ...
	cookieUserIDValueType string
)

// constants ...
const (
	// UserIDCookieName ...
	UserIDCookieName = "user"
	// UserCTXName ...
	UserCTXName cookieUserIDValueType = "user_in_context"
	// UserIDDefaultValue ...
	UserIDDefaultValue = "default_user"
)

// vars ...
var (
	// encryptor ...
	encryptor *Encryptor
	// log ...
	log *zap.Logger
)

// init ...
func init() {
	log, _ = zap.NewProduction()
}

// generateRandom byte slice
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("rand read: %w", err)
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
		return fmt.Errorf("generate key: %w", err)
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("initialize cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return fmt.Errorf("initialize GCM encryptor: %w", err)
	}

	nonce, err := generateRandom(aesGCM.NonceSize())
	if err != nil {
		return fmt.Errorf("initialize GCM nonce: %w", err)
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
		return fmt.Errorf("hex decode: %v", err)
	}

	src, err := e.GCM.Open(nil, e.nonce, dst, nil)
	if err != nil {
		return fmt.Errorf("gcm open: %v", err)
	}

	*to = string(src)
	return nil
}

// AuthMiddleware ...
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rawUserID string

		if err := NewEncryptor(); err != nil {
			log.Error(fmt.Sprintf("new encryptor: %v", err))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCTXName, UserIDDefaultValue)))
			return
		}

		if user, err := r.Cookie(UserIDCookieName); err != nil {
			rawUserID = uuid.New().String()
		} else if err = encryptor.DecodeUUID(user.Value, &rawUserID); err != nil {
			log.Debug(fmt.Sprintf("decode: %v", err))
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
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCTXName, rawUserID)))
	})
}
