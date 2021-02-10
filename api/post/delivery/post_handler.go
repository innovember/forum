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
	commentUcase  post.CommentUsecase
}

func NewPostHandler(postUcase post.PostUsecase, userUcase user.UserUsecase,
	rateUcase post.RateUsecase, categoryUcase post.CategoryUsecase,
	commentUcase post.CommentUsecase) *PostHandler {
	return &PostHandler{
		postUcase:     postUcase,
		userUcase:     userUcase,
		rateUcase:     rateUcase,
		categoryUcase: categoryUcase,
		commentUcase:  commentUcase}
}

func (ph *PostHandler) Configure(mux *http.ServeMux, mw *middleware.MiddlewareManager) {
	// Posts
	mux.HandleFunc("/api/post/create", mw.SetHeaders(mw.AuthorizedOnly(ph.CreatePostHandler)))
	mux.HandleFunc("/api/posts", mw.SetHeaders(ph.GetPostsHandler))
	mux.HandleFunc("/api/post/", mw.SetHeaders(ph.GetPostHandler))
	mux.HandleFunc("/api/post/rate", mw.SetHeaders(mw.AuthorizedOnly(ph.RatePostHandler)))
	mux.HandleFunc("/api/post/filter", mw.SetHeaders(ph.FilterPosts))
	mux.HandleFunc("/api/post/edit", mw.SetHeaders(mw.AuthorizedOnly(ph.EditPostHandler)))
	mux.HandleFunc("/api/post/delete", mw.SetHeaders(mw.AuthorizedOnly(ph.DeletePostHandler)))
	mux.HandleFunc("/api/categories", mw.SetHeaders(ph.GetAllCategoriesHandler))
	// Comments
	mux.HandleFunc("/api/comment/create", mw.SetHeaders(mw.AuthorizedOnly(ph.CreateCommentHandler)))
	mux.HandleFunc("/api/comment/filter", mw.SetHeaders(ph.FilterComments))
	// mux.HandleFunc("/api/comment/edit", mw.SetHeaders(mw.AuthorizedOnly(ph.EditCommentHandler)))
	// mux.HandleFunc("/api/comment/delete", mw.SetHeaders(mw.AuthorizedOnly(ph.DeleteCommentHandler)))
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
		return
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
		EditedAt:   0,
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
		input         models.InputRate
		rating        models.Rating
		err           error
		status        int
		user          *models.User
		cookie        *http.Cookie
		isRatedBefore bool
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if input.Reaction != -1 && input.Reaction != 1 {
		response.Error(w, http.StatusBadRequest, errors.New("only 1 or -1 values accepted"))
		return
	}
	cookie, _ = r.Cookie(config.SessionCookieName)
	if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
		response.Error(w, status, err)
		return
	}
	switch input.Reaction {
	case 1:
		isRatedBefore, err = ph.rateUcase.IsRatedBefore(input.ID, user.ID, input.Reaction)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if isRatedBefore {
			if err = ph.rateUcase.DeleteRateFromPost(input.ID, user.ID, input.Reaction); err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
			response.Success(w, "rate cancelled due to re-voting", http.StatusOK, nil)
			return
		}
	case -1:
		isRatedBefore, err = ph.rateUcase.IsRatedBefore(input.ID, user.ID, input.Reaction)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if isRatedBefore {
			if err = ph.rateUcase.DeleteRateFromPost(input.ID, user.ID, input.Reaction); err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
			response.Success(w, "rate cancelled due to re-voting", http.StatusOK, nil)
			return
		}
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

func (ph *PostHandler) FilterPosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.FilterPostsFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (ph *PostHandler) FilterPostsFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input  models.InputFilterPost
		posts  []models.Post
		status int
		err    error
		cookie *http.Cookie
		user   *models.User
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
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
	switch input.Option {
	case "categories":
		if posts, status, err = ph.postUcase.GetPostsByCategories(input.Categories, user.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "date":
		if posts, status, err = ph.postUcase.GetPostsByDate(input.Date, user.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "rating":
		if posts, status, err = ph.postUcase.GetPostsByRating(input.Rating, user.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "author":
		if posts, status, err = ph.postUcase.GetAllPostsByAuthorID(input.AuthorID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "user":
		if posts, status, err = ph.postUcase.GetRatedPostsByUser(user.ID, input.UserRating); err != nil {
			response.Error(w, status, err)
			return
		}
	default:
		response.Error(w, http.StatusBadRequest, errors.New("option error in filter"))
		return
	}
	response.Success(w, "filtered posts by"+input.Option, status, posts)
	return
}

func (ph *PostHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.CreateCommentHandlerFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (ph *PostHandler) CreateCommentHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input      models.InputComment
		comment    models.Comment
		newComment *models.Comment
		now        = time.Now().Unix()
		status     int
		err        error
		cookie     *http.Cookie
		user       *models.User
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	cookie, _ = r.Cookie(config.SessionCookieName)
	if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
		response.Error(w, status, err)
		return
	}
	comment = models.Comment{
		AuthorID:  user.ID,
		PostID:    input.PostID,
		Content:   input.Content,
		CreatedAt: now,
		EditedAt:  0,
	}
	if newComment, status, err = ph.commentUcase.Create(user.ID, &comment); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "new comment created", status, newComment)
	return
}

func (ph *PostHandler) FilterComments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.FilterCommentsFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (ph *PostHandler) FilterCommentsFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input    models.InputFindComment
		comments []models.Comment
		status   int
		err      error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	switch input.Option {
	case "post":
		if comments, status, err = ph.commentUcase.GetCommentsByPostID(input.PostID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "user":
		if comments, status, err = ph.commentUcase.GetCommentsByAuthorID(input.UserID); err != nil {
			response.Error(w, status, err)
			return
		}
	default:
		response.Error(w, http.StatusBadRequest, errors.New("option error in filter"))
		return
	}
	response.Success(w, "filtered comments by"+input.Option, status, comments)
	return
}

func (ph *PostHandler) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var (
			input      models.InputPost
			post       models.Post
			editedPost *models.Post
			now        = time.Now().Unix()
			status     int
			err        error
			cookie     *http.Cookie
			user       *models.User
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if input.AuthorID != user.ID {
			response.Error(w, http.StatusForbidden, errors.New("can't edit another user's post"))
			return
		}
		post = models.Post{
			ID:       input.ID,
			AuthorID: input.AuthorID,
			Author:   user,
			Title:    input.Title,
			Content:  input.Content,
			EditedAt: now,
		}
		if err = ph.categoryUcase.Update(post.ID, input.Categories); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if editedPost, status, err = ph.postUcase.Update(&post); err != nil {
			response.Error(w, status, err)
			return
		}
		response.Success(w, "post has been edited", status, editedPost)
	} else {
		http.Error(w, "Only PUT method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			input  models.InputPost
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if input.AuthorID != user.ID {
			response.Error(w, http.StatusForbidden, errors.New("can't delete another user's post"))
			return
		}
		if err = ph.categoryUcase.DeleteFromPostCategoriesBridge(input.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if status, err = ph.postUcase.Delete(input.ID); err != nil {
			response.Error(w, status, err)
			return
		}
		response.Success(w, "post has been deleted", status, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}
