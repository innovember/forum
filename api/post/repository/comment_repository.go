package repository

import (
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
)

type CommentDBRepository struct {
	dbConn *sql.DB
}

func NewCommentDBRepository(conn *sql.DB) post.CommentRepository {
	return &CommentDBRepository{dbConn: conn}
}

func (cr *CommentDBRepository) Create(userID int64, comment *models.Comment) (*models.Comment, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = cr.dbConn.Exec(`
	INSERT INTO comments(author_id,post_id,content, created_at,edited_at)
	VALUES(?,?,?,?,?)`, comment.AuthorID, comment.PostID, comment.Content, comment.CreatedAt, comment.EditedAt); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if comment.ID, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return comment, http.StatusCreated, nil
	}
	return nil, http.StatusBadRequest, errors.New("comment hasn't been created")
}

func (cr *CommentDBRepository) GetCommentsByPostID(postID int64) (comments []models.Comment, status int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = cr.dbConn.Query(`
	SELECT * FROM comments 
	WHERE post_id = ?
	ORDER BY created_at DESC`, postID,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Comment
		rows.Scan(&c.ID, &c.AuthorID, &c.PostID, &c.Content, &c.CreatedAt, &c.EditedAt)
		if status, err = cr.GetAuthor(&c); err != nil {
			return nil, status, err
		}
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return comments, http.StatusOK, nil
}

func (cr *CommentDBRepository) GetAuthor(comment *models.Comment) (status int, err error) {
	var (
		user models.User
	)
	if err = cr.dbConn.QueryRow(`
	SELECT id,username,email,created_at,last_active FROM users WHERE id = ?`, comment.AuthorID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastActive); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusUnauthorized, errors.New("cant find author of post")
		}
		return http.StatusInternalServerError, err
	}
	comment.Author = &user
	return http.StatusOK, nil
}

func (cr *CommentDBRepository) GetCommentsByAuthorID(authorID int64) (comments []models.Comment, status int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = cr.dbConn.Query(`
		SELECT *,
		FROM comments
		WHERE author_id = $1
		ORDER BY created_at DESC
		`, authorID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Comment
		rows.Scan(&c.ID, &c.AuthorID, &c.PostID, &c.Content, &c.CreatedAt, &c.EditedAt)
		if status, err = cr.GetAuthor(&c); err != nil {
			return nil, status, err
		}
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return comments, http.StatusOK, nil
}

func (cr *CommentDBRepository) GetCommentsNumberByPostID(postID int64) (commentsNumber int, err error) {
	if err = cr.dbConn.QueryRow(`
	SELECT COUNT(id)
	FROM comments
	WHERE post_id = ?`, postID).Scan(&commentsNumber); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return commentsNumber, nil
}
