package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		stop := time.Now()
		var userID string = UserIDDefaultValue
		ctx := r.Context()
		if ctx != nil {
			if user := ctx.Value(UserCTXName); user != nil {
				userID = user.(string)
			}
		}
		log.Printf("request(path: %v, user: %v, duration: %v)", r.URL, userID, stop.Sub(start))
	})
}
