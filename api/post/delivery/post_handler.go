package delivery

import (
	"encoding/json"
	"errors"
	"github.com/innovember/forum/api/config"
	"github.com/innovember/forum/api/middleware"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"github.com/innovember/forum/api/response"
	"github.com/innovember/forum/api/user"
	"net/http"
	"strconv"
	"time"
)

type PostHandler struct {
	postUcase     post.PostUsecase
	userUcase     user.UserUsecase
	rateUcase     post.RateUsecase
	categoryUcase post.CategoryUsecase
}

func NewPostHandler(postUcase post.PostUsecase, userUcase user.UserUsecase,
	rateUcase post.RateUsecase, categoryUcase post.CategoryUsecase) *PostHandler {
	return &PostHandler{
		postUcase:     postUcase,
		userUcase:     userUcase,
		rateUcase:     rateUcase,
		categoryUcase: categoryUcase}
}

func (ph *PostHandler) Configure(mux *http.ServeMux, mw *middleware.MiddlewareManager) {
	mux.HandleFunc("/api/post/create", mw.SetHeaders(mw.AuthorizedOnly(ph.CreatePostHandler)))
	mux.HandleFunc("/api/posts", mw.SetHeaders(ph.GetPostsHandler))
	mux.HandleFunc("/api/post/", mw.SetHeaders(ph.GetPostHandler))
	mux.HandleFunc("/api/post/rate", mw.SetHeaders(mw.AuthorizedOnly(ph.RatePostHandler)))
	mux.HandleFunc("/api/categories", mw.SetHeaders(ph.GetAllCategoriesHandler))
	// mux.HandleFunc("/api/post/filter", mw.SetHeaders(ph.FilterPosts))
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

func (ph *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			posts  []models.Post
		)
		cookie, err = r.Cookie(config.SessionCookieName)
		if err != nil {
			user = &models.User{ID: -1}
		} else {
			if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
				user = &models.User{ID: -1}
			}
		}
		posts, status, err = ph.postUcase.GetAllPosts(user.ID)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Success(w, "all posts", status, posts)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) RatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.RatePostHandlerFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (ph *PostHandler) RatePostHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input  models.InputRate
		rating models.Rating
		err    error
		status int
		user   *models.User
		cookie *http.Cookie
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}
	if input.Reaction != -1 || input.Reaction != 1 {
		response.Error(w, http.StatusBadRequest, errors.New("only 1 or -1 values accepted"))
		return
	}
	cookie, _ = r.Cookie(config.SessionCookieName)
	if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = ph.rateUcase.RatePost(input.ID, user.ID, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if rating.Rating, rating.UserRating, err = ph.rateUcase.GetRating(input.ID, user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "post has been rated", http.StatusOK, rating)
	return
}

func (ph *PostHandler) GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status     int
			err        error
			categories []models.Category
		)
		categories, status, err = ph.categoryUcase.GetAllCategories()
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Success(w, "all categories", status, categories)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			post   *models.Post
			postID int
		)
		_id := r.URL.Path[len("/api/post/"):]
		if postID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post id doesn't exist"))
			return
		}
		cookie, err = r.Cookie(config.SessionCookieName)
		if err != nil {
			user = &models.User{ID: -1}
		} else {
			if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
				user = &models.User{ID: -1}
			}
		}
		post, status, err = ph.postUcase.GetPostByID(user.ID, int64(postID))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Success(w, "post with id", status, post)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}
