package middleware

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

const UserIDCookieName = "user"

// AuthMiddleware ...
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, err := r.Cookie(UserIDCookieName); err == nil {
			log.Print(user.Value)
		} else if err == http.ErrNoCookie {
			id := uuid.New().String()
			c := &http.Cookie{
				Name:  UserIDCookieName,
				Value: id,

				Path: "/",
			}
			r.AddCookie(c)
		}
		next.ServeHTTP(w, r)
	})
}
