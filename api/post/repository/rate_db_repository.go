package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
)

type RateDBRepository struct {
	dbConn *sql.DB
}

func NewRateDBRepository(conn *sql.DB) post.RateRepository {
	return &RateDBRepository{dbConn: conn}
}

func (rr *RateDBRepository) GetPostRating(postID int64, userID int64) (rating int, userRating int, err error) {
	if err = rr.dbConn.QueryRow(`
	SELECT TOTAL(rate) AS rating,
	IFNULL ((SELECT rate
		 	FROM post_rating
			WHERE user_id = $1 AND
			post_id = $2), 0) AS userRating
	FROM post_rating
	WHERE post_id = $2`,
		userID, postID).Scan(
		&rating, &userRating); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, errors.New("post not found")
		}
		return 0, 0, err
	}
	return rating, userRating, nil

}

func (rr *RateDBRepository) RatePost(postID int64, userID int64, vote int) (int64, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
		rateID       int64
	)
	if result, err = rr.dbConn.Exec(`
		REPLACE INTO post_rating(id, user_id, post_id,rate)
		VALUES(
			(SELECT id FROM post_rating
				WHERE user_id = $1 AND post_id = $2),
			$1,$2,$3)`,
		userID, postID, vote); err != nil {
		return 0, err
	}
	if rateID, err = result.LastInsertId(); err != nil {
		return 0, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return 0, err
	}
	if rowsAffected > 0 {
		return rateID, nil
	}
	return 0, errors.New("cant set new rate for post")
}

func (rr *RateDBRepository) IsRatedBefore(postID int64, userID int64, vote int) (bool, error) {
	var (
		err  error
		rate int
	)
	if err = rr.dbConn.QueryRow(`SELECT rate FROM post_rating
								 WHERE post_id = ?
								 AND user_id = ?
								 AND rate = ?`, postID, userID, vote).Scan(
		&rate); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if rate == vote {
		return true, nil
	}
	return false, nil
}

func (rr *RateDBRepository) DeleteRateFromPost(postID int64, userID int64, vote int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = rr.dbConn.Exec(
		`DELETE FROM post_rating
		WHERE post_id = ?
		AND user_id = ?
		AND rate = ?`, postID, userID, vote,
	); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return errors.New("post with such rate not found")
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not delete the rate")
}

func (rr *RateDBRepository) GetPostRatingByID(rateID int64) (postRating *models.PostRating, status int, err error) {
	var (
		pr models.PostRating
	)
	if err = rr.dbConn.QueryRow(`
		SELECT *
		FROM post_rating
		WHERE id = ?
		`, rateID).Scan(&pr.ID, &pr.UserID,
		&pr.PostID, &pr.Rate); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("cant find rating")
		}
		return nil, http.StatusInternalServerError, err
	}
	if status, err = rr.GetAuthor(&pr); err != nil {
		return nil, status, err
	}
	return &pr, http.StatusOK, nil
}

func (rr *RateDBRepository) GetAuthor(postRating *models.PostRating) (status int, err error) {
	var (
		user models.User
	)
	if err = rr.dbConn.QueryRow(`
	SELECT id,username,email,created_at,last_active FROM users WHERE id = ?`, postRating.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastActive); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("cant find author of rating")
		}
		return http.StatusInternalServerError, err
	}
	postRating.Author = &user
	return http.StatusOK, nil
}

func (rr *RateDBRepository) DeleteRatesByPostID(postID int64) (status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`DELETE FROM post_rating
								WHERE post_id = ?`,
		postID); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("rates not found")
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
	return http.StatusNotModified, errors.New("could not delete the rates")
}
