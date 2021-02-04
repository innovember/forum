package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
	"strings"
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
	)
	if result, err = pr.dbConn.Exec(`
	INSERT INTO posts(author_id,title, content, created_at)
	VALUES(?,?,?,?)`, post.AuthorID, post.Title, post.Content, now); err != nil {
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

func (pr *PostDBRepository) GetAllPosts(userID int64) (posts []models.Post, status int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = pr.dbConn.Query(`
		SELECT *,
		(SELECT TOTAL(rate)
			FROM post_rating
			WHERE post_id = posts.id) AS rating,
		IFNULL (
				(
					SELECT rate
					FROM post_rating
					WHERE post_id = posts.id
					AND user_id = $1
					),0) AS userRating
		FROM posts
		ORDER BY created_at DESC
		`, userID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.PostRating, &p.UserRating)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return posts, http.StatusOK, nil
}

func (pr *PostDBRepository) GetAuthor(post *models.Post) (status int, err error) {
	var (
		user models.User
	)
	if err = pr.dbConn.QueryRow(`
	SELECT id,username,email,created_at,last_active FROM users WHERE id = ?`, post.AuthorID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastActive); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusUnauthorized, errors.New("cant find author of post")
		}
		return http.StatusInternalServerError, err
	}
	post.Author = &user
	return http.StatusOK, nil
}

func (pr *PostDBRepository) GetCategories(post *models.Post) (status int, err error) {
	var (
		rows       *sql.Rows
		categories []models.Category
	)
	if rows, err = pr.dbConn.Query(`
		SELECT c.id,c.name
		FROM categories c
		LEFT JOIN posts_categories_bridge pcb
		ON pcb.post_id = ?
		WHERE c.id = pcb.category_id`,
		post.ID); err != nil {
		return http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}
	err = rows.Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	post.Categories = categories
	return http.StatusOK, nil
}

func (pr *PostDBRepository) GetPostByID(userID int64, postID int64) (post *models.Post, status int, err error) {
	var (
		p        models.Post
		rateRepo = NewRateDBRepository(pr.dbConn)
	)
	if err = pr.dbConn.QueryRow(`
	SELECT * FROM posts WHERE id = ?`, postID,
	).Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusUnauthorized, errors.New("post not found")
		}
		return nil, http.StatusInternalServerError, err
	}
	if status, err = pr.GetAuthor(&p); err != nil {
		return nil, status, err
	}
	if p.PostRating, p.UserRating, err = rateRepo.GetPostRating(postID, userID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &p, http.StatusOK, nil
}

func (pr *PostDBRepository) GetPostsByCategories(categories []string, userID int64) (posts []models.Post, status int, err error) {
	var (
		rows           *sql.Rows
		categoriesList string = fmt.Sprintf("\"%s\"", strings.Join(categories, "\", \""))
		rateRepo              = NewRateDBRepository(pr.dbConn)
	)
	query := `
		SELECT p.*
		FROM posts_categories_bridge as pcb
		INNER JOIN posts as p
		ON p.id = pcb.post_id
		INNER JOIN categories as c
		ON c.id=pcb.category_id
		WHERE c.name in (?)
		GROUP BY p.id
		HAVING COUNT(DISTINCT c.id) = ?
		`
	if rows, err = pr.dbConn.Query(query, categoriesList, len(categories)); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.PostRating, p.UserRating, err = rateRepo.GetPostRating(p.ID, userID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return posts, http.StatusOK, nil
}

func (pr *PostDBRepository) GetPostsByRating(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = pr.dbConn.Query(`
		SELECT *,
		(SELECT TOTAL(rate)
			FROM post_rating
			WHERE post_id = posts.id) AS rating,
		IFNULL (
				(
					SELECT rate
					FROM post_rating
					WHERE post_id = posts.id
					AND user_id = $1
					),0) AS userRating
		FROM posts
		ORDER BY rating $2
		`, userID, orderBy); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.PostRating, &p.UserRating)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return posts, http.StatusOK, nil
}

func (pr *PostDBRepository) GetPostsByDate(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = pr.dbConn.Query(`
		SELECT *,
		(SELECT TOTAL(rate)
			FROM post_rating
			WHERE post_id = posts.id) AS rating,
		IFNULL (
				(
					SELECT rate
					FROM post_rating
					WHERE post_id = posts.id
					AND user_id = $1
					),0) AS userRating
		FROM posts
		ORDER BY created_at $2
		`, userID, orderBy); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.PostRating, &p.UserRating)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return posts, http.StatusOK, nil
}
