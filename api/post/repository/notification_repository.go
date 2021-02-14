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
	INSERT INTO notifications(receiver_id, post_id,rate_id,comment_id,created_at)
	VALUES(?,?,?,?,?)`, notification.ReceiverID, notification.PostID,
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
