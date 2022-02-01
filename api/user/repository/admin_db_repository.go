package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
)

type AdminDBRepository struct {
	dbConn *sql.DB
}

func NewAdminDBRepository(conn *sql.DB) user.AdminRepository {
	return &AdminDBRepository{dbConn: conn}
}

func (ar *AdminDBRepository) UpgradeRole(requestID int64) (err error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		userID int64
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if err = tx.QueryRow(`SELECT user_id
						 FROM role_requests
						 WHERE id = ?
		`, requestID).Scan(&userID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`UPDATE users
						 SET role = 1
						 WHERE id = ? 
		`, userID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM role_requests
						 WHERE id = ?
		`, requestID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *AdminDBRepository) GetAllRoleRequests() (roleRequests []models.RoleRequest, err error) {
	var (
		ctx      context.Context
		tx       *sql.Tx
		rows     *sql.Rows
		userRepo = NewUserDBRepository(ar.dbConn)
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT *
							 FROM role_requests
		`); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r models.RoleRequest
		err = rows.Scan(&r.ID, &r.UserID, &r.CreatedAt, &r.Pending)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		r.User, err = userRepo.GetUserByID(r.UserID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		roleRequests = append(roleRequests, r)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return roleRequests, nil
}

func (ar *AdminDBRepository) DeleteRoleRequest(requestID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM role_requests
						 WHERE id = ?
		`, requestID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *AdminDBRepository) GetAllPostReports() (postReports []models.PostReport, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT r.*, p.title
							 FROM post_reports AS r
							 JOIN posts AS p
							 ON p.id = r.post_id
		`); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pr models.PostReport
		err = rows.Scan(&pr.ID, &pr.ModeratorID,
			&pr.PostID, &pr.Pending,
			&pr.PostTitle)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		postReports = append(postReports, pr)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return postReports, nil
}

func (ar *AdminDBRepository) AcceptPostReport(postReportID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM posts
						 WHERE id IN (
							 SELECT post_id
							 FROM post_reports
							 WHERE id = ?
						 )
		`, postReportID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM post_reports
						 WHERE id = ?
		`, postReportID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *AdminDBRepository) DismissPostReport(postReportID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM post_reports
						 WHERE id = ?
		`, postReportID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *AdminDBRepository) GetAllModerators() (moderators []models.User, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT id, username,email,created_at, last_active,role
							 FROM users
							 WHERE role = 1`); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email,
			&u.CreatedAt, &u.LastActive,
			&u.Role)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		moderators = append(moderators, u)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return moderators, tx.Commit()
}

func (ar *AdminDBRepository) DemoteModerator(moderatorID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE users
						 SET role = 0
						 WHERE id = ? 
		`, moderatorID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
