package delivery

import (
	"encoding/json"
	"errors"
	"github.com/innovember/forum/api/config"
	"github.com/innovember/forum/api/middleware"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"github.com/innovember/forum/api/response"
	"github.com/innovember/forum/api/security"
	"github.com/innovember/forum/api/user"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserHandler struct {
	userUcase         user.UserUsecase
	adminUcase        user.AdminUsecase
	moderatorUcase    user.ModeratorUsecase
	postUcase         post.PostUsecase
	rateUcase         post.RateUsecase
	categoryUcase     post.CategoryUsecase
	commentUcase      post.CommentUsecase
	notificationUcase post.NotificationUsecase
	commentRateUcase  post.RateCommentUsecase
}

func NewUserHandler(
	userUcase user.UserUsecase,
	adminUcase user.AdminUsecase,
	moderatorUcase user.ModeratorUsecase,
	postUcase post.PostUsecase,
	rateUcase post.RateUsecase,
	categoryUcase post.CategoryUsecase,
	commentUcase post.CommentUsecase,
	notificationUcase post.NotificationUsecase,
	commentRateUcase post.RateCommentUsecase) *UserHandler {
	return &UserHandler{
		userUcase:         userUcase,
		adminUcase:        adminUcase,
		moderatorUcase:    moderatorUcase,
		postUcase:         postUcase,
		rateUcase:         rateUcase,
		categoryUcase:     categoryUcase,
		commentUcase:      commentUcase,
		notificationUcase: notificationUcase,
		commentRateUcase:  commentRateUcase,
	}
}

func (uh *UserHandler) Configure(mux *http.ServeMux, mw *middleware.MiddlewareManager) {
	// auth
	mux.HandleFunc("/api/auth/signup", mw.SetHeaders(uh.CreateUserHandler))
	mux.HandleFunc("/api/auth/signin", mw.SetHeaders(uh.SignIn))
	mux.HandleFunc("/api/auth/signout", mw.SetHeaders(mw.AuthorizedOnly(uh.SignOut)))
	mux.HandleFunc("/api/auth/me", mw.SetHeaders(mw.AuthorizedOnly(uh.Me)))
	// user's info
	mux.HandleFunc("/api/users", mw.SetHeaders(uh.GetAllUsers))
	mux.HandleFunc("/api/user/", mw.SetHeaders(uh.GetUserByID))
	// user's role
	mux.HandleFunc("/api/request/add", mw.SetHeaders(mw.AuthorizedOnly(uh.CreateRoleRequest)))
	mux.HandleFunc("/api/request/delete", mw.SetHeaders(mw.AuthorizedOnly(uh.DeleteRoleRequest)))
	mux.HandleFunc("/api/request", mw.SetHeaders(mw.AuthorizedOnly(uh.GetRoleRequest)))
	// admin
	mux.HandleFunc("/api/admin/requests", mw.SetHeaders(mw.AuthorizedOnly(uh.GetRoleRequests)))
	mux.HandleFunc("/api/admin/request/dismiss/", mw.SetHeaders(mw.AuthorizedOnly(uh.DismissRoleRequest)))
	mux.HandleFunc("/api/admin/request/accept/", mw.SetHeaders(mw.AuthorizedOnly(uh.AcceptRoleRequest)))

	mux.HandleFunc("/api/admin/post/reports", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllPostReports)))
	mux.HandleFunc("/api/admin/post/report/dismiss/", mw.SetHeaders(mw.AuthorizedOnly(uh.DismissPostReport)))
	mux.HandleFunc("/api/admin/post/report/accept/", mw.SetHeaders(mw.AuthorizedOnly(uh.AcceptPostReport)))

	mux.HandleFunc("/api/admin/post/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeletePostByAdmin)))
	mux.HandleFunc("/api/admin/comment/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeleteCommentByAdmin)))

	mux.HandleFunc("/api/admin/moderators", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllModerators)))
	mux.HandleFunc("/api/admin/demote/moderator/", mw.SetHeaders(mw.AuthorizedOnly(uh.DemoteModerator)))

	mux.HandleFunc("/api/admin/categories", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllCategories)))
	mux.HandleFunc("/api/admin/category/add", mw.SetHeaders(mw.AuthorizedOnly(uh.CreateNewCategory)))
	mux.HandleFunc("/api/admin/category/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeleteCategory)))

	// moderator
	mux.HandleFunc("/api/moderator/reports", mw.SetHeaders(mw.AuthorizedOnly(uh.MyReports)))
	mux.HandleFunc("/api/moderator/report/post/create", mw.SetHeaders(mw.AuthorizedOnly(uh.CreatePostReport)))
	mux.HandleFunc("/api/moderator/report/post/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeletePostReport)))
	mux.HandleFunc("/api/moderator/post/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeletePostByModerator)))

	// moderator -> post reviewing
	mux.HandleFunc("/api/moderator/posts/unapproved", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllUnapprovedPosts)))
	mux.HandleFunc("/api/moderator/post/approve/", mw.SetHeaders(mw.AuthorizedOnly(uh.ApprovePost)))
	mux.HandleFunc("/api/moderator/post/ban/", mw.SetHeaders(mw.AuthorizedOnly(uh.BanPost)))

	// user notifications
	// mux.HandleFunc("/api/user/notifications/admin", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllNotificationsFromAdmin)))
	// mux.HandleFunc("/api/user/notifications/admin/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeleteNotificationsFromAdmin)))
	// mux.HandleFunc("/api/user/notifications/moderator", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllNotificationsFromModerator)))
	// mux.HandleFunc("/api/user/notifications/moderator/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeleteNotificationsFromModerator)))
	// mux.HandleFunc("/api/moderator/notifications/admin", mw.SetHeaders(mw.AuthorizedOnly(uh.GetAllNotificationsForModeratorFromAdmin)))
	// mux.HandleFunc("/api/moderator/notifications/admin/delete/", mw.SetHeaders(mw.AuthorizedOnly(uh.DeleteNotificationsForModeratorFromAdmin)))
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
		Role:      config.RoleUser,
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

func (uh *UserHandler) CreateRoleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = uh.userUcase.CreateRoleRequest(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "new role request created", http.StatusCreated, nil)
	} else {
		http.Error(w, "Only POST method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DeleteRoleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = uh.userUcase.DeleteRoleRequest(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "role request has been removed", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) GetRoleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status      int
			err         error
			cookie      *http.Cookie
			user        *models.User
			roleRequest *models.RoleRequest
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if roleRequest, err = uh.userUcase.GetRoleRequestByUserID(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "get role request by user id", http.StatusOK, roleRequest)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) GetRoleRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status       int
			err          error
			cookie       *http.Cookie
			user         *models.User
			roleRequests []models.RoleRequest
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		if roleRequests, err = uh.adminUcase.GetAllRoleRequests(); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "get role requests", http.StatusOK, roleRequests)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DismissRoleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status        int
			err           error
			cookie        *http.Cookie
			user          *models.User
			roleRequestID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/admin/request/dismiss/"):]
		if roleRequestID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("request id doesn't exist"))
			return
		}
		if err = uh.adminUcase.DeleteRoleRequest(int64(roleRequestID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "role request has been removed", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) AcceptRoleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var (
			status        int
			err           error
			cookie        *http.Cookie
			user          *models.User
			roleRequestID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/admin/request/accept/"):]
		if roleRequestID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("request id doesn't exist"))
			return
		}
		if err = uh.adminUcase.UpgradeRole(int64(roleRequestID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "role request has been accepted", http.StatusOK, nil)
	} else {
		http.Error(w, "Only PUT method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DeletePostByAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			postID int
			post   *models.Post
		)
		_id := r.URL.Path[len("/api/admin/post/delete/"):]
		if postID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post id doesn't exist"))
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		post, status, err = uh.postUcase.GetPostByID(user.ID, int64(postID))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if err = uh.categoryUcase.DeleteFromPostCategoriesBridge(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.rateUcase.DeleteRatesByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.commentUcase.DeleteCommentsByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.notificationUcase.DeleteNotificationsByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.commentRateUcase.DeleteCommentsRateByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if status, err = uh.postUcase.Delete(post.ID); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post has been deleted", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DeleteCommentByAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status    int
			err       error
			cookie    *http.Cookie
			user      *models.User
			commentID int
			comment   *models.Comment
		)
		_id := r.URL.Path[len("/api/admin/comment/delete/"):]
		if commentID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("comment id doesn't exist"))
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		if comment, status, err = uh.commentUcase.GetCommentByID(user.ID, int64(commentID)); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = uh.notificationUcase.DeleteNotificationsByCommentID(comment.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.commentRateUcase.DeleteRatesByCommentID(comment.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.commentUcase.Delete(comment.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "comment has been deleted", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DeletePostByModerator(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			postID int
			post   *models.Post
		)
		_id := r.URL.Path[len("/api/moderator/post/delete/"):]
		if postID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post id doesn't exist"))
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		post, status, err = uh.postUcase.GetPostByID(user.ID, int64(postID))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if err = uh.categoryUcase.DeleteFromPostCategoriesBridge(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.rateUcase.DeleteRatesByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.commentUcase.DeleteCommentsByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.notificationUcase.DeleteNotificationsByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.commentRateUcase.DeleteCommentsRateByPostID(post.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if status, err = uh.postUcase.Delete(post.ID); err != nil {
			response.Error(w, status, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post has been deleted", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) CreatePostReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var (
			input  models.PostReport
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		if err = uh.moderatorUcase.CreatePostReport(&input); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post report created", http.StatusCreated, input)
	} else {
		http.Error(w, "Only POST method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DeletePostReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status       int
			err          error
			cookie       *http.Cookie
			user         *models.User
			postReportID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/moderator/report/post/delete/"):]
		if postReportID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post report id doesn't exist"))
			return
		}
		if err = uh.moderatorUcase.DeletePostReport(int64(postReportID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post report has been removed", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) MyReports(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status      int
			err         error
			cookie      *http.Cookie
			user        *models.User
			postReports []models.PostReport
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		if postReports, err = uh.moderatorUcase.GetMyReports(user.ID); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "fetched all post reports by moderator id", http.StatusOK, postReports)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status     int
			err        error
			cookie     *http.Cookie
			user       *models.User
			categories []models.Category
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		if categories, status, err = uh.categoryUcase.GetAllCategories(); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "all categories", status, categories)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var (
			input  models.Category
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		if err = uh.categoryUcase.CreateNewCategory(input.Name); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "new category has been created", http.StatusCreated, input)
	} else {
		http.Error(w, "Only POST method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status     int
			err        error
			cookie     *http.Cookie
			user       *models.User
			categoryID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/admin/category/delete/"):]
		if categoryID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("category id doesn't exist"))
			return
		}
		if err = uh.categoryUcase.DeleteCategoryByID(int64(categoryID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "category has been removed", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) GetAllPostReports(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status      int
			err         error
			cookie      *http.Cookie
			user        *models.User
			postReports []models.PostReport
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		if postReports, err = uh.adminUcase.GetAllPostReports(); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "get all post reports", http.StatusOK, postReports)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DismissPostReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var (
			status       int
			err          error
			cookie       *http.Cookie
			user         *models.User
			postReportID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/admin/post/report/dismiss/"):]
		if postReportID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post report id doesn't exist"))
			return
		}
		if err = uh.adminUcase.DismissPostReport(int64(postReportID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post report has been removed", http.StatusOK, nil)
	} else {
		http.Error(w, "Only DELETE method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) AcceptPostReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var (
			status       int
			err          error
			cookie       *http.Cookie
			user         *models.User
			postReportID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/admin/request/accept/"):]
		if postReportID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("request id doesn't exist"))
			return
		}
		if err = uh.adminUcase.AcceptPostReport(int64(postReportID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post report has been accepted", http.StatusOK, nil)
	} else {
		http.Error(w, "Only PUT method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) GetAllModerators(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status     int
			err        error
			cookie     *http.Cookie
			user       *models.User
			moderators []models.User
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		if moderators, err = uh.adminUcase.GetAllModerators(); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "get all moderators", http.StatusOK, moderators)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) DemoteModerator(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var (
			status      int
			err         error
			cookie      *http.Cookie
			user        *models.User
			moderatorID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only admin users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/admin/demote/moderator/"):]
		if moderatorID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post report id doesn't exist"))
			return
		}
		if err = uh.adminUcase.DemoteModerator(int64(moderatorID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "moderator has been demoted", http.StatusOK, nil)
	} else {
		http.Error(w, "Only PUT method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) GetAllUnapprovedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			posts  []models.Post
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		if posts, err = uh.moderatorUcase.GetAllUnapprovedPosts(); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "all unapproved posts", http.StatusOK, posts)
	} else {
		http.Error(w, "Only GET method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) ApprovePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			postID int
		)
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		_id := r.URL.Path[len("/api/moderator/post/approve/"):]
		if postID, err = strconv.Atoi(_id); err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("post id doesn't exist"))
			return
		}
		if err = uh.moderatorUcase.ApprovePost(int64(postID)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post has been approved", http.StatusOK, nil)
	} else {
		http.Error(w, "Only PUT method allowed, return to main page", 405)
		return
	}
}

func (uh *UserHandler) BanPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var (
			status int
			err    error
			cookie *http.Cookie
			user   *models.User
			postID int
			input  models.InputPost
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, _ = r.Cookie(config.SessionCookieName)
		if user, status, err = uh.userUcase.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		if user.Role != config.RoleModerator {
			response.Error(w, http.StatusForbidden, errors.New("not enough privileges,only moderator users allowed"))
			return
		}
		if err = uh.moderatorUcase.BanPost(int64(postID), input.Bans); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if err = uh.userUcase.UpdateActivity(user.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "post has been banned", http.StatusOK, nil)
	} else {
		http.Error(w, "Only POST method allowed, return to main page", 405)
		return
	}
}
