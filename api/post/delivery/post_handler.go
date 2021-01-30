package delivery

import (
	"encoding/json"
	"github.com/innovember/forum/api/config"
	"github.com/innovember/forum/api/middleware"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"github.com/innovember/forum/api/response"
	"github.com/innovember/forum/api/user"
	"net/http"
	"time"
)

type PostHandler struct {
	postUcase post.PostUsecase
	userUcase user.UserUsecase
}

func NewPostHandler(postUcase post.PostUsecase, userUcase user.UserUsecase) *PostHandler {
	return &PostHandler{
		postUcase: postUcase,
		userUcase: userUcase}
}

func (ph *PostHandler) Configure(mux *http.ServeMux, mw *middleware.MiddlewareManager) {
	mux.HandleFunc("/api/post/create", mw.SetHeaders(mw.AuthorizedOnly(ph.CreatePostHandler)))
}

func (ph *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.CreatePostHandlerFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (ph *PostHandler) CreatePostHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input   models.InputPost
		post    models.Post
		newPost *models.Post
		now     = time.Now().Unix()
		status  int
		err     error
		cookie  *http.Cookie
		user    *models.User
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}
	cookie, _ = r.Cookie(config.SessionCookieName)
	if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
		response.Error(w, status, err)
		return
	}
	post = models.Post{
		AuthorID:   user.ID,
		Author:     user,
		Title:      input.Title,
		Content:    input.Content,
		CreatedAt:  now,
		PostRating: 0,
	}
	if newPost, status, err = ph.postUcase.Create(&post, input.Categories); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "new post created", status, newPost)
	return
}
