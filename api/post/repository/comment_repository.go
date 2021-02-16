package repository

import (
	"context"
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
			return http.StatusNotFound, errors.New("cant find author of post")
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
		SELECT *
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

func (cr *CommentDBRepository) Update(comment *models.Comment) (editedComment *models.Comment, status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`UPDATE comments
							SET content = ?,
							edited_at = ?
							WHERE post_id = ?
							and id = ?`,
		comment.Content, comment.EditedAt, comment.PostID, comment.ID); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, http.StatusInternalServerError, errors.New("comment not found")
		}
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		if err := tx.Commit(); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return comment, http.StatusOK, nil
	}
	return nil, http.StatusNotModified, errors.New("could not update the comment")
}

func (cr *CommentDBRepository) GetCommentByID(commentID int64) (comment *models.Comment, status int, err error) {
	var c models.Comment
	if err = cr.dbConn.QueryRow(`
	SELECT * FROM comments WHERE id = ?`, commentID,
	).Scan(&c.ID, &c.AuthorID, &c.PostID, &c.Content, &c.CreatedAt, &c.EditedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("comment not found")
		}
		return nil, http.StatusInternalServerError, err
	}
	if status, err = cr.GetAuthor(&c); err != nil {
		return nil, status, err
	}
	return &c, http.StatusOK, nil
}

func (cr *CommentDBRepository) Delete(commentID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM comments
								WHERE id = ?`,
		commentID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *CommentDBRepository) DeleteCommentByPostID(postID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM comments
								WHERE post_id = ?`,
		postID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
