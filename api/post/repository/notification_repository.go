package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
	"time"
)

type NotificationDBRepository struct {
	dbConn *sql.DB
}

func NewNotificationDBRepository(conn *sql.DB) post.NotificationRepository {
	return &NotificationDBRepository{dbConn: conn}
}

func (nr *NotificationDBRepository) Create(notification *models.Notification) (*models.Notification, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		now          = time.Now().Unix()
		err          error
	)
	if result, err = nr.dbConn.Exec(`
	INSERT INTO notifications(receiver_id, post_id,
		rate_id,comment_id,comment_rate_id,created_at)
	VALUES(?,?,?,?,?,?)`, notification.ReceiverID, notification.PostID,
		notification.RateID, notification.CommentID,
		notification.CommentRateID, now); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if notification.ID, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return notification, http.StatusCreated, nil
	}
	return nil, http.StatusBadRequest, errors.New("notification hasn't been created")
}

func (nr *NotificationDBRepository) DeleteAllNotifications(receiverID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications
								WHERE receiver_id = ?`,
		receiverID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (nr *NotificationDBRepository) GetAllNotifications(receiverID int64) (notifications []models.Notification, status int, err error) {
	var (
		rows            *sql.Rows
		postRepo        = NewPostDBRepository(nr.dbConn)
		commentRepo     = NewCommentDBRepository(nr.dbConn)
		rateRepo        = NewRateDBRepository(nr.dbConn)
		commentRateRepo = NewRateCommentDBRepository(nr.dbConn)
	)
	if rows, err = nr.dbConn.Query(`
		SELECT *
		FROM notifications
		WHERE receiver_id = ?
		ORDER BY created_at DESC
		`, receiverID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var n models.Notification
		rows.Scan(&n.ID, &n.ReceiverID, &n.PostID, &n.RateID,
			&n.CommentID, &n.CommentRateID, &n.CreatedAt)
		if n.Post, status, err = postRepo.GetPostByID(receiverID, n.PostID); err != nil {
			return nil, status, err
		}
		if n.RateID != 0 {
			if n.PostRating, status, err = rateRepo.GetPostRatingByID(n.RateID); err != nil {
				return nil, status, err
			}
		}
		if n.CommentID != 0 {
			if n.Comment, status, err = commentRepo.GetCommentByID(receiverID, n.CommentID); err != nil {
				return nil, status, err
			}
		}
		if n.CommentRateID != 0 {
			if n.CommentRating, status, err = commentRateRepo.GetCommentRatingByID(n.CommentRateID); err != nil {
				return nil, status, err
			}
		}
		notifications = append(notifications, n)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return notifications, http.StatusOK, nil
}

func (nr *NotificationDBRepository) DeleteNotificationsByPostID(postID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications
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

func (nr *NotificationDBRepository) DeleteNotificationsByRateID(rateID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications
								WHERE rate_id = ?`,
		rateID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (nr *NotificationDBRepository) DeleteNotificationsByCommentID(commentID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications
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

func (nr *NotificationDBRepository) DeleteNotificationsByCommentRateID(commentRateID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications
								WHERE comment_rate_id = ?`,
		commentRateID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
