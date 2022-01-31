package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
)

type UserNotificationDBRepository struct {
	dbConn *sql.DB
}

func NewUserNotificationDBRepository(conn *sql.DB) user.UserNotificationRepository {
	return &UserNotificationDBRepository{dbConn: conn}
}

func (ur *UserNotificationDBRepository) CreateRoleNotification(roleNotification *models.RoleNotification) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		now = time.Now().Unix()
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO notifications_roles(receiver_id, accepted,
		declined,demoted,created_at)
	VALUES(?,?,?,?,?)`, roleNotification.ReceiverID, roleNotification.Accepted,
		roleNotification.Declined, roleNotification.Demoted, now); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserNotificationDBRepository) CreatePostReportNotification(postReportNotification *models.PostReportNotification) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		now = time.Now().Unix()
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO notifications_reports(receiver_id, approved,
		deleted,created_at)
	VALUES(?,?,?,?)`, postReportNotification.ReceiverID, postReportNotification.Approved,
		postReportNotification.Deleted, now); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserNotificationDBRepository) DeleteAllRoleNotifications(userID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications_roles
						 WHERE receiver_id = ?
		`, userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserNotificationDBRepository) DeleteAllPostReportNotifications(userID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications_reports
						 WHERE receiver_id = ?
		`, userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserNotificationDBRepository) GetRoleNotifications(userID int64) (roleNotifications []models.RoleNotification, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT *
							 FROM notifications_roles
							 WHERE receiver_id = ?`,
		userID); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rn models.RoleNotification
		err = rows.Scan(&rn.ID, &rn.ReceiverID, &rn.Accepted,
			&rn.Declined, &rn.Demoted, &rn.CreatedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		roleNotifications = append(roleNotifications, rn)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return roleNotifications, tx.Commit()
}

func (ur *UserNotificationDBRepository) GetPostReportNotifications(userID int64) (postReportNotifications []models.PostReportNotification, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT *
							 FROM notifications_reports
							 WHERE receiver_id = ?`,
		userID); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pn models.PostReportNotification
		err = rows.Scan(&pn.ID, &pn.ReceiverID, &pn.Approved,
			&pn.Deleted, &pn.CreatedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		postReportNotifications = append(postReportNotifications, pn)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return postReportNotifications, tx.Commit()
}

func (ur *UserNotificationDBRepository) CreatePostNotification(postNotification *models.PostNotification) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		now = time.Now().Unix()
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO notifications_posts(receiver_id, approved,
		deleted,created_at)
	VALUES(?,?,?,?,?)`,
		postNotification.ReceiverID,
		postNotification.Approved,
		postNotification.Banned,
		postNotification.Deleted,
		now); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserNotificationDBRepository) DeleteAllPostNotifications(userID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM notifications_posts
						 WHERE receiver_id = ?
		`, userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserNotificationDBRepository) GetPostNotifications(userID int64) (postNotifications []models.PostNotification, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT *
							 FROM notifications_posts
							 WHERE receiver_id = ?`,
		userID); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pn models.PostNotification
		err = rows.Scan(&pn.ID, &pn.ReceiverID, &pn.Approved,
			&pn.Banned, &pn.Deleted, &pn.CreatedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		postNotifications = append(postNotifications, pn)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return postNotifications, tx.Commit()
}
