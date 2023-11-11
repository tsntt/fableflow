package middleware

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func Cors(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return next(w, r)
		}

		h := w.Header()

		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")

		// CORS preflight.
		if r.Method == http.MethodOptions {
			h.Set("Access-Control-Allow-Methods", "GET,PATCH,POST,DELETE,HEAD")
			h.Set("Access-Control-Allow-Headers", "authorization,content-type")
			h.Set("Access-Control-Max-Age", "86400")
			return nil
		}

		return next(w, r)
	}
}
