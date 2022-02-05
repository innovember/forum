package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
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
	INSERT INTO posts(author_id,title, content, created_at,edited_at, is_image,image_path,is_approved)
	VALUES(?,?,?,?,?,?,?,?)`, post.AuthorID, post.Title,
		post.Content, now, post.EditedAt,
		post.IsImage, post.ImagePath, post.IsApproved); err != nil {
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
		rows        *sql.Rows
		commentRepo = NewCommentDBRepository(pr.dbConn)
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
		WHERE is_approved = 1
		ORDER BY created_at DESC
		`, userID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned,
			&p.PostRating, &p.UserRating)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
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
			return http.StatusNotFound, errors.New("cant find author of post")
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
		p           models.Post
		rateRepo    = NewRateDBRepository(pr.dbConn)
		commentRepo = NewCommentDBRepository(pr.dbConn)
	)
	if err = pr.dbConn.QueryRow(`
	SELECT * FROM posts WHERE id = ?
	AND is_approved = 1`, postID,
	).Scan(&p.ID, &p.AuthorID, &p.Title,
		&p.Content, &p.CreatedAt,
		&p.EditedAt, &p.IsImage, &p.ImagePath,
		&p.IsApproved, &p.IsBanned); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("post not found")
		}
		return nil, http.StatusInternalServerError, err
	}
	if status, err = pr.GetAuthor(&p); err != nil {
		return nil, status, err
	}
	if p.PostRating, p.UserRating, err = rateRepo.GetPostRating(postID, userID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if status, err = pr.GetCategories(&p); err != nil {
		return nil, status, err
	}
	if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
		return nil, status, err
	}
	return &p, http.StatusOK, nil
}

func (pr *PostDBRepository) GetPostsByCategories(categories []string, userID int64) (posts []models.Post, status int, err error) {
	var (
		rows           *sql.Rows
		categoriesList = fmt.Sprintf("\"%s\"", strings.Join(categories, "\", \""))
		rateRepo       = NewRateDBRepository(pr.dbConn)
		commentRepo    = NewCommentDBRepository(pr.dbConn)
	)
	query := fmt.Sprintf(`
		SELECT p.*
		FROM posts_categories_bridge as pcb
		INNER JOIN posts as p
		ON p.id = pcb.post_id
		INNER JOIN categories as c
		ON c.id=pcb.category_id
		WHERE c.name in (%s)
		AND p.is_approved = 1
		GROUP BY p.id
		HAVING COUNT(DISTINCT c.id) = %d
		ORDER BY p.created_at DESC`, categoriesList, len(categories))
	if rows, err = pr.dbConn.Query(query); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.PostRating, p.UserRating, err = rateRepo.GetPostRating(p.ID, userID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
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

func (pr *PostDBRepository) GetPostsByRating(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	var (
		rows        *sql.Rows
		commentRepo = NewCommentDBRepository(pr.dbConn)
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
		WHERE is_approved = 1
		ORDER BY rating $2
		`, userID, orderBy); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned,
			&p.PostRating, &p.UserRating)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
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
		rows        *sql.Rows
		commentRepo = NewCommentDBRepository(pr.dbConn)
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
		WHERE is_approved = 1
		ORDER BY created_at $2
		`, userID, orderBy); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned,
			&p.PostRating, &p.UserRating)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
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

func (pr *PostDBRepository) GetAllPostsByAuthorID(authorID int64, userID int64) (posts []models.Post, status int, err error) {
	var (
		rows        *sql.Rows
		commentRepo = NewCommentDBRepository(pr.dbConn)
		rateRepo    = NewRateDBRepository(pr.dbConn)
	)
	if rows, err = pr.dbConn.Query(`
		SELECT *
		FROM posts
		WHERE author_id = ?
		AND is_approved = 1
		ORDER BY created_at DESC
		`, authorID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.PostRating, p.UserRating, err = rateRepo.GetPostRating(p.ID, userID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
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

func (pr *PostDBRepository) GetRatedPostsByUser(userID int64, orderBy string, requestorID int64) (posts []models.Post, status int, err error) {
	var (
		rows        *sql.Rows
		vote        int
		rateRepo    = NewRateDBRepository(pr.dbConn)
		commentRepo = NewCommentDBRepository(pr.dbConn)
	)
	if orderBy == "upvoted" {
		vote = 1
	} else if orderBy == "downvoted" {
		vote = -1
	}
	query := fmt.Sprintf(`
		SELECT p.* FROM posts AS p
		INNER JOIN post_rating AS pr ON p.id = pr.post_id
		WHERE pr.user_id = %d AND pr.rate = %d
		AND p.is_approved = 1
		ORDER BY p.created_at DESC
		`, userID, vote)
	if rows, err = pr.dbConn.Query(query); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned)
		if status, err = pr.GetAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = pr.GetCategories(&p); err != nil {
			return nil, status, err
		}
		if p.PostRating, p.UserRating, err = rateRepo.GetPostRating(p.ID, requestorID); err != nil {
			return nil, status, err
		}
		if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
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

func (pr *PostDBRepository) Update(post *models.Post) (editedPost *models.Post, status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`UPDATE posts
							SET title = ?,
							content = ?,
							edited_at = ?,
							is_image = ?,
							image_path = ?
							WHERE id = ?`,
		post.Title, post.Content, post.EditedAt,
		post.IsImage, post.ImagePath, post.ID); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, http.StatusInternalServerError, errors.New("post not found")
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
		return post, http.StatusOK, nil
	}
	return nil, http.StatusNotModified, errors.New("could not update the post")
}

func (pr *PostDBRepository) Delete(postID int64) (status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`DELETE FROM posts
								WHERE id = ?`,
		postID); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("post not found")
		}
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, nil
	}
	if rowsAffected > 0 {
		if err := tx.Commit(); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the post")
}

func (pr *PostDBRepository) GetBannedPostsByCategories(categories []string) (posts []models.Post, status int, err error) {
	var (
		rows           *sql.Rows
		categoriesList string = fmt.Sprintf("\"%s\"", strings.Join(categories, "\", \""))
	)
	query := fmt.Sprintf(`
		SELECT p.*
		FROM posts_bans_bridge as pbb
		INNER JOIN posts as p
		ON p.id = pbb.post_id
		INNER JOIN bans as b
		ON b.id=pbb.ban_id
		WHERE b.name in (%s)
		AND p.is_approved = 0
		AND p.is_banned = 1
		GROUP BY p.id
		HAVING COUNT(DISTINCT b.id) = %d
		ORDER BY p.created_at DESC`, categoriesList, len(categories))
	if rows, err = pr.dbConn.Query(query); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved, &p.IsBanned)
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

func (pr *PostDBRepository) DeletePostReportByPostID(postID int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM post_reports
						 WHERE post_id = ?
		`, postID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
