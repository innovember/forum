package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
)

type RateCommentDBRepository struct {
	dbConn *sql.DB
}

func NewRateCommentDBRepository(conn *sql.DB) post.RateCommentRepository {
	return &RateCommentDBRepository{dbConn: conn}
}

func (rr *RateCommentDBRepository) GetCommentRating(commentID int64, userID int64) (rating int, userRating int, err error) {
	if err = rr.dbConn.QueryRow(`
	SELECT TOTAL(rate) AS rating,
	IFNULL ((SELECT rate
		 	FROM comment_rating
			WHERE user_id = $1 AND
			comment_id = $2), 0) AS userRating
	FROM comment_rating
	WHERE comment_id = $2`,
		userID, commentID).Scan(
		&rating, &userRating); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, errors.New("comment not found")
		}
		return 0, 0, err
	}
	return rating, userRating, nil

}

func (rr *RateCommentDBRepository) RateComment(commentID int64, userID int64, vote int, postID int64) (int64, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
		rateID       int64
	)
	if result, err = rr.dbConn.Exec(`
		REPLACE INTO comment_rating(id, user_id,post_id, comment_id,rate)
		VALUES(
			(SELECT id FROM comment_rating
				WHERE user_id = $1 AND comment_id = $2
			AND post_id = $4),
			$1,$4,$2,$3)`,
		userID, commentID, vote, postID); err != nil {
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
	return 0, errors.New("cant set new rate for comment")
}

func (rr *RateCommentDBRepository) IsRatedBefore(commentID int64, userID int64, vote int) (bool, error) {
	var (
		err  error
		rate int
	)
	if err = rr.dbConn.QueryRow(`SELECT rate FROM comment_rating
								 WHERE comment_id = ?
								 AND user_id = ?
								 AND rate = ?`, commentID, userID, vote).Scan(
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

func (rr *RateCommentDBRepository) DeleteRateFromComment(commentID int64, userID int64, vote int) (err error) {
	var (
		ctx           context.Context
		tx            *sql.Tx
		commentRateID int64
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if err = tx.QueryRow(`SELECT id
							FROM comment_rating
							WHERE comment_id = ?
		AND user_id = ?
		AND rate = ?`, commentID, userID, vote).Scan(&commentRateID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications
    WHERE comment_rate_id = ?`, commentRateID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM comment_rating
		WHERE comment_id = ?
		AND user_id = ?
		AND rate = ?`, commentID, userID, vote); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *RateCommentDBRepository) GetCommentRatingByID(commentRateID int64) (commentRating *models.CommentRating, status int, err error) {
	var (
		cr models.CommentRating
	)
	if err = rr.dbConn.QueryRow(`
		SELECT *
		FROM comment_rating
		WHERE id = ?
		`, commentRateID).Scan(&cr.ID, &cr.UserID, &cr.PostID,
		&cr.CommentID, &cr.Rate); err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("cant find comment rating")
		}
		return nil, http.StatusInternalServerError, err
	}
	if status, err = rr.GetAuthor(&cr); err != nil {
		return nil, status, err
	}
	return &cr, http.StatusOK, nil
}

func (rr *RateCommentDBRepository) GetAuthor(commentRating *models.CommentRating) (status int, err error) {
	var (
		user models.User
	)
	if err = rr.dbConn.QueryRow(`
	SELECT id,username,email,created_at,last_active FROM users WHERE id = ?`, commentRating.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastActive); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("cant find author of rating")
		}
		return http.StatusInternalServerError, err
	}
	commentRating.Author = &user
	return http.StatusOK, nil
}

func (rr *RateCommentDBRepository) DeleteRatesByCommentID(commentID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM comment_rating
								WHERE comment_id = ?`,
		commentID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *RateCommentDBRepository) DeleteCommentsRateByPostID(postID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM comment_rating
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
