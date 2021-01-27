package middleware

import (
	"github.com/innovember/forum/api/config"
	"net/http"
)

func (mw *MiddlewareManager) SetHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", config.ClientURL)
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next(w, r)
	}
}

// func (mw *MiddlewareManager) SetHeaders(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if origin := r.Header.Get("Origin"); origin != "" {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
// 			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 			// w.Header().Set("Access-Control-Allow-Origin", config.ClientURL)
// 			if r.Method == "OPTIONS" {
// 				w.WriteHeader(200)
// 				return
// 			}
// 			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 			next(w, r)
// 		}
// 	}
// }
