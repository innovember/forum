package repository

import (
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
	"net/http"
	"time"
)

type UserDBRepository struct {
	dbConn *sql.DB
}

func NewUserDBRepository(conn *sql.DB) user.UserRepository {
	return &UserDBRepository{dbConn: conn}
}

func (ur *UserDBRepository) Create(user *models.User) (status int, err error) {
	var (
		result       sql.Result
		rowsAffected int64
		now          int64 = time.Now().Unix()
	)
	if result, err = ur.dbConn.Exec(`
	INSERT INTO users (
			username,
			password,
			email,
			created_at,
			last_active,
			session_id
		)
	VALUES (?, ?, ?, ?, ?,?)`,
		user.Username, user.Password, user.Email, now, now, user.SessionID,
	); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusBadRequest, errors.New("cant create new user")
}

func (ur *UserDBRepository) CheckByUsernameOrEmail(user *models.User) (status int, err error) {
	username := ur.dbConn.QueryRow("SELECT username FROM users WHERE username = ?", user.Username).Scan(&user.Username)
	email := ur.dbConn.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&user.Email)
	if username == nil && email != nil {
		return http.StatusConflict, errors.New("username already exist")
	} else if email == nil && username != nil {
		return http.StatusConflict, errors.New("email already exist")
	} else if username == nil && email == nil {
		return http.StatusConflict, errors.New("username and email already exist")
	}
	return http.StatusOK, nil
}
