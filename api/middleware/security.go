package middleware

import (
	"errors"
	"github.com/innovember/forum/api/config"
	// "github.com/innovember/forum/api/db"
	// "github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/response"
	"net/http"
)

func (mw *MiddlewareManager) AuthorizedOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(config.SessionCookieName)
		if err != nil {
			response.Error(w, http.StatusForbidden, errors.New("user not authorized"))
		}
		next(w, r)
	}
}
