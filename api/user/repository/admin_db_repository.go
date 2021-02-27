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

func (ar *AdminDBRepository) UpgradeRole(userID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE users
						 SET role = 1
						 WHERE id = ? 
		`, userID); err != nil {
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
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
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
			return nil, err
		}
		roleRequests = append(roleRequests, r)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return roleRequests, nil
}

func (ar *AdminDBRepository) DeleteRoleRequest(userID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ar.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM role_requests
						 WHERE user_id = ?
		`, userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
