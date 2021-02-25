package delivery

import (
	"encoding/json"
	"errors"
	"github.com/innovember/forum/api/config"
	"github.com/innovember/forum/api/middleware"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/response"
	"github.com/innovember/forum/api/security"
	"github.com/innovember/forum/api/user"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserHandler struct {
	userUcase user.UserUsecase
}

func NewUserHandler(userUcase user.UserUsecase) *UserHandler {
	return &UserHandler{userUcase: userUcase}
}

func (uh *UserHandler) Configure(mux *http.ServeMux, mw *middleware.MiddlewareManager) {
	mux.HandleFunc("/api/auth/signup", mw.SetHeaders(uh.CreateUserHandler))
	mux.HandleFunc("/api/users", mw.SetHeaders(uh.GetAllUsers))
	mux.HandleFunc("/api/auth/signin", mw.SetHeaders(uh.SignIn))
	mux.HandleFunc("/api/auth/signout", mw.SetHeaders(mw.AuthorizedOnly(uh.SignOut)))
	mux.HandleFunc("/api/auth/me", mw.SetHeaders(mw.AuthorizedOnly(uh.Me)))
	mux.HandleFunc("/api/user/", mw.SetHeaders(uh.GetUserByID))
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uh.CreateUserHandlerFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (uh *UserHandler) CreateUserHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input          models.InputUserSignUp
		hashedPassword string
		status         int
		err            error
		adminAuthToken string = os.Getenv("ADMIN_AUTH_TOKEN")
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if hashedPassword, err = security.Hash(input.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{
		Username:  input.Username,
		Password:  hashedPassword,
		Email:     input.Email,
		SessionID: "",
		Role:      config.RoleGuest,
	}
	if adminAuthToken != "" && adminAuthToken == input.AdminAuthToken {
		user.Role = config.RoleAdmin
	}
	if status, err = uh.userUcase.Create(&user); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	response.Success(w, "new user has been created", status, user)
	return
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		users, status, err := uh.userUcase.GetAllUsers()
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Success(w, "all users", status, users)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uh.SignInFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (uh *UserHandler) SignInFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input        models.InputUserSignIn
		user         *models.User
		userPassword string
		cookie       string
		newUUID      string
		err          error
		status       int
		expiresAt    = time.Now().Add(config.SessionExpiration).Unix()
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if user, status, err = uh.userUcase.FindUserByUsername(input.Username); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = uh.userUcase.CheckSessionByUsername(user.Username); err != nil {
		response.Error(w, status, err)
		return
	}
	if userPassword, status, err = uh.userUcase.GetPassword(user.Username); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = security.VerifyPassword(userPassword, input.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	cookie, newUUID = security.GenerateCookie(r.Cookie(config.SessionCookieName))
	if err = uh.userUcase.UpdateSession(user.ID, newUUID, expiresAt); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, "user logged in", http.StatusOK, user)
}

func (uh *UserHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.SignOutFunc(w, r)
	default:
		http.Error(w, "Only GET method allowed, return to main page", 405)
	}
}
func (uh *UserHandler) SignOutFunc(w http.ResponseWriter, r *http.Request) {
	var (
		user      *models.User
		err       error
		status    int
		cookie    *http.Cookie
		expiresAt = time.Now().Add(config.SessionExpiration).Unix()
	)
	if cookie, err = r.Cookie(config.SessionCookieName); err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
		response.Error(w, status, err)
	}
	if err = uh.userUcase.UpdateSession(user.ID, "", expiresAt); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	cookie = &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	response.Success(w, "user logged out", http.StatusOK, nil)
	return
}

func (uh *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.MeFunc(w, r)
	default:
		http.Error(w, "Only GET method allowed, return to main page", 405)
	}
}

func (uh *UserHandler) MeFunc(w http.ResponseWriter, r *http.Request) {
	var (
		user   *models.User
		status int
		err    error
		cookie *http.Cookie
	)
	cookie, _ = r.Cookie(config.SessionCookieName)
	if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "get user info successfully", status, user)
}

func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			user   *models.User
			err    error
			userID int
		)
		_id := r.URL.Path[len("/api/user/"):]
		if userID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("user id doesn't exist"))
			return
		}
		user, err = uh.userUcase.GetUserByID(int64(userID))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Success(w, "fetch user by id", http.StatusOK, user)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}
