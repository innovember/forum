package repository

import (
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/post"
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

func (rr *RateDBRepository) RatePost(postID int64, userID int64, vote int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = rr.dbConn.Exec(`
		REPLACE INTO post_rating(id, user_id, post_id,rate)
		VALUES(
			(SELECT id FROM post_rating
				WHERE user_id = $1 AND post_id = $2),
			$1,$2,$3)`,
		userID, postID, vote); err != nil {
		return err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("cant set new rate for post")
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
