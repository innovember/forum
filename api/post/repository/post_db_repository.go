package repository

import (
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
	"time"
)

type PostDBRepository struct {
	dbConn *sql.DB
}

func NewPostDBRepository(conn *sql.DB) post.PostRepository {
	return &PostDBRepository{dbConn: conn}
}

func (pr *PostDBRepository) Create(post *models.Post, categories []string) (*models.Post, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		now          = time.Now().Unix()
		categoryRepo = NewCategoryDBRepository(pr.dbConn)
		err          error
		status       int
	)
	if result, err = pr.dbConn.Exec(`
	INSERT INTO posts(author_id,title, content, created_at)
	VALUES(?,?,?,?)`, post.AuthorID, post.Title, post.Content, post.CreatedAt); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if post.ID, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if err = categoryRepo.Create(post.ID, categories); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return post, http.StatusCreated, nil
	}
	return nil, http.StatusBadRequest, errors.New("post hasn't been created")
}
