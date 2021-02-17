package session

import (
	"context"
	"database/sql"
	"log"
	// "github.com/innovember/forum/api/config"
	userRepo "github.com/innovember/forum/api/user/repository"
	"time"
)

func Init(dbConn *sql.DB) {
	go CheckSessionExpiration(dbConn)
}

func ResetAll(dbConn *sql.DB) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE users
						 SET session_id = ?,
						 expires_at = ?
						 WHERE id IN
						 (
							 IFNULL(
							 (SELECT id
						  	 FROM users)
						  	 ,0)
						  )`, "", 0); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func CheckSessionExpiration(dbConn *sql.DB) {
	var (
		ctx            context.Context
		tx             *sql.Tx
		rows           *sql.Rows
		err            error
		userID         int64
		sessionValue   string
		expiresAt      int64
		userRepository = userRepo.NewUserDBRepository(dbConn)
	)
	for {
		time.Sleep(time.Minute * 5)
		now := time.Now().Unix()
		ctx = context.Background()
		if tx, err = dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
			log.Println(err)
			return
		}
		if rows, err = tx.Query(`SELECT id,session_id, expires_at
								FROM users
								WHERE expires_at < ?`,
			now); err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&userID, &sessionValue, &expiresAt)
			if err = userRepository.UpdateSession(userID, sessionValue, expiresAt); err != nil {
				log.Println(err)
				return
			}
		}
		if err = tx.Commit(); err != nil {
			log.Println(err)
			return
		}
	}
}
