package middleware

import (
	"errors"
	"github.com/innovember/forum/api/config"
	// "github.com/innovember/forum/api/db"
	// "github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/db"
	"github.com/innovember/forum/api/response"
	userRepo "github.com/innovember/forum/api/user/repository"
	userUsecase "github.com/innovember/forum/api/user/usecases"
	"net/http"
)

func (mw *MiddlewareManager) AuthorizedOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			cookie *http.Cookie
		)
		//Repository
		userRepository := userRepo.NewUserDBRepository(db.DBConn)

		//Usecases
		userUcase := userUsecase.NewUserUsecase(userRepository)

		cookie, err = r.Cookie(config.SessionCookieName)
		if err != nil {
			response.Error(w, http.StatusForbidden, errors.New("session not found, user not authorized"))
			return
		}
		if cookie.Value == "" {
			response.Error(w, http.StatusForbidden, errors.New("user not authorized"))
			return
		}
		if _, _, err = userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, http.StatusForbidden, errors.New("session not valid,user not authorized"))
			return
		}
		next(w, r)
	}
}
