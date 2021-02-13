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
		return http.StatusCreated, nil
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

func (ur *UserDBRepository) GetAllUsers() (users []models.User, err error) {
	var rows *sql.Rows
	if rows, err = ur.dbConn.Query(`SELECT id, username,email,created_at, last_active FROM users`); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.LastActive)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserDBRepository) GetPassword(username string) (password string, status int, err error) {
	if err = ur.dbConn.QueryRow(`
	SELECT password FROM users WHERE username = ?`, username).Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			return "", http.StatusNotFound, errors.New("password not found for username:" + username)
		}
		return "", http.StatusInternalServerError, err
	}
	return password, http.StatusOK, nil
}

func (ur *UserDBRepository) FindUserByUsername(username string) (*models.User, int, error) {
	var (
		user models.User
		err  error
	)
	if err = ur.dbConn.QueryRow(`
	SELECT id,username,email FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("user not found for username:" + username)
		}
		return nil, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

func (ur *UserDBRepository) UpdateSession(userID int64, sessionValue string) (err error) {
	var (
		result       sql.Result
		rowsAffected int64
		now          int64 = time.Now().Unix()
	)
	if result, err = ur.dbConn.Exec(`
	UPDATE users SET session_id = ?,last_active = ? WHERE id = ?`, sessionValue, now, userID); err != nil {
		return err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("cant update session")
}

func (ur *UserDBRepository) ValidateSession(sessionValue string) (*models.User, int, error) {
	var (
		user models.User
		err  error
	)
	if err = ur.dbConn.QueryRow(`
	SELECT id,username,email,created_at,last_active FROM users WHERE session_id = ?`, sessionValue).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("user not authorized")
		}
		return nil, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

func (ur *UserDBRepository) CheckSessionByUsername(username string) (status int, err error) {
	var user models.User
	if err = ur.dbConn.QueryRow(`
	SELECT session_id FROM users WHERE username = ?`, username).Scan(&user.SessionID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusInternalServerError, err
	}
	if user.SessionID == "" {
		return http.StatusOK, nil
	}
	return http.StatusForbidden, errors.New("user already authorized")
}

func (ur *UserDBRepository) GetUserByID(userID int64) (*models.User, error) {
	var (
		user models.User
		err  error
	)
	if err = ur.dbConn.QueryRow(`
	SELECT id,username,email,created_at,last_active
	FROM users WHERE id = ?`, userID).Scan(&user.ID, &user.Username,
		&user.Email, &user.CreatedAt, &user.LastActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("cant find user with such id")
		}
		return nil, err
	}
	return &user, nil
}
