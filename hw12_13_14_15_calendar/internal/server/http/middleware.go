package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf(
			"%s [%s] %s %s %s \"%s\"",
			r.RemoteAddr,
			time.Now().String(),
			r.Method,
			r.URL.Path,
			r.Proto,
			r.Header.Get("User-Agent"),
		)
		fmt.Println(s)

		next.ServeHTTP(w, r)
	})
}
