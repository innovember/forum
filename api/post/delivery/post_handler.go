package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/innovember/forum/api/config"
	"github.com/innovember/forum/api/middleware"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"github.com/innovember/forum/api/response"
	"github.com/innovember/forum/api/user"
	uuid "github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PostHandler struct {
	postUcase         post.PostUsecase
	userUcase         user.UserUsecase
	rateUcase         post.RateUsecase
	categoryUcase     post.CategoryUsecase
	commentUcase      post.CommentUsecase
	notificationUcase post.NotificationUsecase
	commentRateUcase  post.RateCommentUsecase
}

func NewPostHandler(postUcase post.PostUsecase, userUcase user.UserUsecase,
	rateUcase post.RateUsecase, categoryUcase post.CategoryUsecase,
	commentUcase post.CommentUsecase, notificationUcase post.NotificationUsecase,
	commentRateUcase post.RateCommentUsecase) *PostHandler {
	return &PostHandler{
		postUcase:         postUcase,
		userUcase:         userUcase,
		rateUcase:         rateUcase,
		categoryUcase:     categoryUcase,
		commentUcase:      commentUcase,
		notificationUcase: notificationUcase,
		commentRateUcase:  commentRateUcase,
	}
}

func (ph *PostHandler) Configure(mux *http.ServeMux, mw *middleware.MiddlewareManager) {
	// Posts
	mux.HandleFunc("/api/post/create", mw.SetHeaders(mw.AuthorizedOnly(ph.CreatePostHandler)))
	mux.HandleFunc("/api/posts", mw.SetHeaders(ph.GetPostsHandler))
	mux.HandleFunc("/api/post/", mw.SetHeaders(ph.GetPostHandler))
	mux.HandleFunc("/api/post/rate", mw.SetHeaders(mw.AuthorizedOnly(ph.RatePostHandler)))
	mux.HandleFunc("/api/post/filter", mw.SetHeaders(ph.FilterPosts))
	mux.HandleFunc("/api/post/edit", mw.SetHeaders(mw.AuthorizedOnly(ph.EditPostHandler)))
	mux.HandleFunc("/api/post/delete/", mw.SetHeaders(mw.AuthorizedOnly(ph.DeletePostHandler)))
	mux.HandleFunc("/api/categories", mw.SetHeaders(ph.GetAllCategoriesHandler))
	// Comments
	mux.HandleFunc("/api/comment/create", mw.SetHeaders(mw.AuthorizedOnly(ph.CreateCommentHandler)))
	mux.HandleFunc("/api/comment/filter", mw.SetHeaders(ph.FilterComments))
	mux.HandleFunc("/api/comment/edit", mw.SetHeaders(mw.AuthorizedOnly(ph.EditCommentHandler)))
	mux.HandleFunc("/api/comment/delete/", mw.SetHeaders(mw.AuthorizedOnly(ph.DeleteCommentHandler)))
	mux.HandleFunc("/api/comment/rate", mw.SetHeaders(mw.AuthorizedOnly(ph.RateCommentHandler)))

	// Notifications
	mux.HandleFunc("/api/notifications", mw.SetHeaders(mw.AuthorizedOnly(ph.GetAllNotificationsHandler)))
	mux.HandleFunc("/api/notifications/delete/", mw.SetHeaders(mw.AuthorizedOnly(ph.DeleteNotificationsHandler)))
	// Images
	mux.HandleFunc("/api/image/upload", mw.SetHeaders(mw.AuthorizedOnly(ph.UploadImageHandler)))
	mux.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./images"))))
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
		IsImage:    input.IsImage,
		ImagePath:  input.ImagePath,
	}
	if newPost, status, err = ph.postUcase.Create(&post, input.Categories); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
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
		rateID        int64
		post          *models.Post
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
	if rateID, err = ph.rateUcase.RatePost(input.ID, user.ID, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if rating.Rating, rating.UserRating, err = ph.rateUcase.GetRating(input.ID, user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err = ph.notificationUcase.DeleteNotificationsByRateID(rateID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if post, status, err = ph.postUcase.GetPostByID(user.ID, input.ID); err != nil {
		response.Error(w, status, err)
		return
	}
	if user.ID != post.AuthorID {
		notification := models.Notification{
			PostID:        input.ID,
			RateID:        rateID,
			CommentID:     0,
			CommentRateID: 0,
		}
		if _, status, err = ph.notificationUcase.Create(&notification); err != nil {
			response.Error(w, status, err)
			return
		}
	}
	if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
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
		if posts, status, err = ph.postUcase.GetAllPostsByAuthorID(input.AuthorID, user.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "user":
		if posts, status, err = ph.postUcase.GetRatedPostsByUser(input.UserID, input.UserRating, user.ID); err != nil {
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
		post       *models.Post
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
	if post, status, err = ph.postUcase.GetPostByID(user.ID, newComment.PostID); err != nil {
		response.Error(w, status, err)
		return
	}
	if user.ID != post.AuthorID {
		notification := models.Notification{
			PostID:        newComment.PostID,
			RateID:        0,
			CommentID:     newComment.ID,
			CommentRateID: 0,
		}
		if _, status, err = ph.notificationUcase.Create(&notification); err != nil {
			response.Error(w, status, err)
			return
		}
	}
	if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
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
		user     *models.User
		cookie   *http.Cookie
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
	case "post":
		if comments, status, err = ph.commentUcase.GetCommentsByPostID(user.ID, input.PostID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "user":
		if comments, status, err = ph.commentUcase.GetCommentsByAuthorID(user.ID, input.UserID); err != nil {
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
			ID:        input.ID,
			AuthorID:  input.AuthorID,
			Author:    user,
			Title:     input.Title,
			Content:   input.Content,
			EditedAt:  now,
			IsImage:   input.IsImage,
			ImagePath: input.ImagePath,
		}
		if err = ph.categoryUcase.Update(post.ID, input.Categories); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if editedPost, status, err = ph.postUcase.Update(&post); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
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
			postID int
			post   *models.Post
		)
		_id := r.URL.Path[len("/api/post/delete/"):]
		if postID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post id doesn't exist"))
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		post, status, err = ph.postUcase.GetPostByID(user.ID, int64(postID))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if post.AuthorID != user.ID {
			response.Error(w, http.StatusForbidden, errors.New("can't delete another user's post"))
			return
		}
		if err = ph.categoryUcase.DeleteFromPostCategoriesBridge(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = ph.rateUcase.DeleteRatesByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = ph.commentUcase.DeleteCommentByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = ph.notificationUcase.DeleteNotificationsByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if status, err = ph.postUcase.Delete(post.ID); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post has been deleted", status, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var (
			input         models.InputComment
			comment       models.Comment
			editedComment *models.Comment
			now           = time.Now().Unix()
			status        int
			err           error
			cookie        *http.Cookie
			user          *models.User
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
			response.Error(w, http.StatusForbidden, errors.New("can't edit another user's comment"))
			return
		}
		comment = models.Comment{
			ID:       input.ID,
			AuthorID: input.AuthorID,
			PostID:   input.PostID,
			Author:   user,
			Content:  input.Content,
			EditedAt: now,
		}
		if editedComment, status, err = ph.commentUcase.Update(&comment); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "comment has been edited", status, editedComment)
	} else {
		http.Error(w, "Only PUT method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status    int
			err       error
			cookie    *http.Cookie
			user      *models.User
			commentID int
			comment   *models.Comment
		)
		_id := r.URL.Path[len("/api/comment/delete/"):]
		if commentID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("comment id doesn't exist"))
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if comment, status, err = ph.commentUcase.GetCommentByID(user.ID, int64(commentID)); err != nil {
			response.Error(w, status, err)
			return
		}
		if comment.AuthorID != user.ID {
			response.Error(w, http.StatusForbidden, errors.New("can't delete another user's comment"))
			return
		}
		if err = ph.notificationUcase.DeleteNotificationsByCommentID(comment.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = ph.commentUcase.Delete(comment.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "comment has been deleted", status, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) GetAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status        int
			err           error
			cookie        *http.Cookie
			user          *models.User
			notifications []models.Notification
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		notifications, status, err = ph.notificationUcase.GetAllNotifications(user.ID)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Success(w, "all notifications", status, notifications)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) DeleteNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = ph.notificationUcase.DeleteAllNotifications(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "notifications has been deleted", status, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var (
			status       int
			err          error
			cookie       *http.Cookie
			user         *models.User
			maxImageSize int64 = config.MaxImageSize
			image        multipart.File
			fileHeader   *multipart.FileHeader
			file         *os.File
			regex        = regexp.MustCompile(`^.*\.(jpg|JPG|jpeg|JPEG|gif|GIF|png|PNG|svg|SVG)$`)
			fileName     string
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = ph.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if r.ContentLength > maxImageSize {
			response.Error(w, http.StatusExpectationFailed, errors.New("image too heavy,limit size to 20MB"))
			return
		}
		if image, fileHeader, err = r.FormFile("image"); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if fileHeader != nil {
			defer image.Close()
		}
		if fileHeader != nil && !regex.MatchString(fileHeader.Filename) {
			response.Error(w, http.StatusUnprocessableEntity, errors.New("invalid file type"))
			return
		}
		fileNameArr := strings.Split(fileHeader.Filename, ".")
		fileExtension := fileNameArr[len(fileNameArr)-1]
		fileName = fmt.Sprint(uuid.NewV4())
		if fileHeader != nil {
			file, err = os.Create(fmt.Sprintf("./images/%s.%s", fileName, fileExtension))
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
			defer file.Close()
			_, err = io.Copy(file, image)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
		}
		if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "image uploaded", http.StatusCreated, fmt.Sprintf("%s/images/%s.%s", config.APIURLDev, fileName, fileExtension))
	} else {
		http.Error(w, "Only POST method allowed, return to main page", 405)
		return
	}
}

func (ph *PostHandler) RateCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.RateCommentHandlerFunc(w, r)
	default:
		http.Error(w, "Only POST method allowed, return to main page", 405)
	}
}

func (ph *PostHandler) RateCommentHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var (
		input         models.InputRate
		rating        models.Rating
		err           error
		status        int
		user          *models.User
		cookie        *http.Cookie
		isRatedBefore bool
		commentRateID int64
		comment       *models.Comment
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
		isRatedBefore, err = ph.commentRateUcase.IsRatedBefore(input.ID, user.ID, input.Reaction)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if isRatedBefore {
			if err = ph.commentRateUcase.DeleteRateFromComment(input.ID, user.ID, input.Reaction); err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
			response.Success(w, "rate cancelled due to re-voting", http.StatusOK, nil)
			return
		}
	case -1:
		isRatedBefore, err = ph.commentRateUcase.IsRatedBefore(input.ID, user.ID, input.Reaction)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if isRatedBefore {
			if err = ph.commentRateUcase.DeleteRateFromComment(input.ID, user.ID, input.Reaction); err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
			response.Success(w, "rate cancelled due to re-voting", http.StatusOK, nil)
			return
		}
	}
	if commentRateID, err = ph.commentRateUcase.RateComment(input.ID, user.ID, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if rating.Rating, rating.UserRating, err = ph.commentRateUcase.GetCommentRating(input.ID, user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err = ph.notificationUcase.DeleteNotificationsByCommentRateID(commentRateID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if comment, status, err = ph.commentUcase.GetCommentByID(user.ID, input.ID); err != nil {
		response.Error(w, status, err)
		return
	}
	if user.ID != comment.AuthorID {
		notification := models.Notification{
			PostID:        input.ID,
			CommentRateID: commentRateID,
			CommentID:     comment.ID,
			RateID:        0,
		}
		if _, status, err = ph.notificationUcase.Create(&notification); err != nil {
			response.Error(w, status, err)
			return
		}
	}
	if err = ph.userUcase.UpdateActivity(user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "comment has been rated", http.StatusOK, rating)
	return
}
