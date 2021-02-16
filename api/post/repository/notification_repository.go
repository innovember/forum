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
		postRepo     = NewPostDBRepository(nr.dbConn)
		post         *models.Post
		status       int
	)
	if post, status, err = postRepo.GetPostByID(-1, notification.PostID); err != nil {
		return nil, status, err
	}
	if result, err = nr.dbConn.Exec(`
	INSERT INTO notifications(receiver_id, post_id,rate_id,comment_id,created_at)
	VALUES(?,?,?,?,?)`, post.AuthorID, notification.PostID,
		notification.RateID, notification.CommentID, now); err != nil {
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

func (nr *NotificationDBRepository) DeleteAllNotifications(receiverID int64) (status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`DELETE FROM notifications
								WHERE receiver_id = ?`,
		receiverID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("notifications not found")
		}
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, nil
	}
	if rowsAffected > 0 {
		if err := tx.Commit(); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the notifications")
}

func (nr *NotificationDBRepository) GetAllNotifications(receiverID int64) (notifications []models.Notification, status int, err error) {
	var (
		rows        *sql.Rows
		postRepo    = NewPostDBRepository(nr.dbConn)
		commentRepo = NewCommentDBRepository(nr.dbConn)
		rateRepo    = NewRateDBRepository(nr.dbConn)
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
			&n.CommentID, &n.CreatedAt)
		if n.Post, status, err = postRepo.GetPostByID(receiverID, n.PostID); err != nil {
			return nil, status, err
		}
		if n.RateID != 0 {
			if n.PostRating, status, err = rateRepo.GetPostRatingByID(n.RateID); err != nil {
				return nil, status, err
			}
		}
		if n.CommentID != 0 {
			if n.Comment, status, err = commentRepo.GetCommentByID(n.CommentID); err != nil {
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

func (nr *NotificationDBRepository) DeleteNotificationsByPostID(postID int64) (status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`DELETE FROM notifications
								WHERE post_id = ?`,
		postID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("notifications not found")
		}
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, nil
	}
	if rowsAffected > 0 {
		if err := tx.Commit(); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the notifications")
}

func (nr *NotificationDBRepository) DeleteNotificationsByRateID(rateID int64) (status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`DELETE FROM notifications
								WHERE rate_id = ?`,
		rateID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("notifications not found")
		}
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, nil
	}
	if rowsAffected > 0 {
		if err := tx.Commit(); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the notifications")
}

func (nr *NotificationDBRepository) DeleteNotificationsByCommentID(commentID int64) (status int, err error) {
	var (
		ctx          context.Context
		tx           *sql.Tx
		result       sql.Result
		rowsAffected int64
	)
	ctx = context.Background()
	if tx, err = nr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return http.StatusInternalServerError, err
	}
	if result, err = tx.Exec(`DELETE FROM notifications
								WHERE comment_id = ?`,
		commentID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("notifications not found")
		}
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, nil
	}
	if rowsAffected > 0 {
		if err := tx.Commit(); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the notifications")
}
