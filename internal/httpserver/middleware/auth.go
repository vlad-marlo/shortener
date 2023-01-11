package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/pkg/encryptor"
)

// types ...
type (
	// cookieUserIDValueType ...
	userCtxKey struct{}
)

// constants ...
const (
	// UserIDCookieName ...
	UserIDCookieName = "user"
	// UserIDDefaultValue ...
	UserIDDefaultValue = "default_user"
)

// vars ...
var (
	// log ...
	log *zap.Logger
)

// init ...
func init() {
	log, _ = zap.NewProduction()
}

// AuthMiddleware ...
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rawUserID string

		if user, err := r.Cookie(UserIDCookieName); err != nil {
			rawUserID = uuid.New().String()
		} else if err = encryptor.Get().DecodeUUID(user.Value, &rawUserID); err != nil {
			log.Debug(fmt.Sprintf("decode: %v", err))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtxKey{}, UserIDDefaultValue)))
			return
		}

		encoded := encryptor.Get().EncodeUUID(rawUserID)
		c := &http.Cookie{
			Name:  UserIDCookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, c)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtxKey{}, rawUserID)))
	})
}
