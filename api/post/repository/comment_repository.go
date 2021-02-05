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
	INSERT INTO comments(author_id,post_id,content, created_at)
	VALUES(?,?,?,?)`, comment.AuthorID, comment.PostID, comment.Content, comment.CreatedAt); err != nil {
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
		rows.Scan(&c.ID, &c.AuthorID, &c.PostID, &c.Content, &c.CreatedAt)
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return comments, http.StatusOK, nil
}
