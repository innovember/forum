package repository

import (
	"context"
	"database/sql"
	"github.com/innovember/forum/api/models"
	postRepo "github.com/innovember/forum/api/post/repository"
	"github.com/innovember/forum/api/user"
	"time"
)

type ModeratorDBRepository struct {
	dbConn *sql.DB
}

func NewModeratorDBRepository(conn *sql.DB) user.ModeratorRepository {
	return &ModeratorDBRepository{dbConn: conn}
}

func (mr *ModeratorDBRepository) CreatePostReport(postReport *models.PostReport) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		now int64 = time.Now().Unix()
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO post_reports(
						moderator_id,post_id,created_at, pending)
						VALUES (
							?,?,?,?
						)
		`, postReport.ModeratorID, postReport.PostID, now, 1); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (mr *ModeratorDBRepository) DeletePostReport(postReportID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
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

func (mr *ModeratorDBRepository) GetMyReports(moderatorID int64) (postReports []models.PostReport, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT *
							 FROM post_reports
							 WHERE moderator_id = ?
		`, moderatorID); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pr models.PostReport
		err = rows.Scan(&pr.ID, &pr.ModeratorID, &pr.PostID, &pr.CreatedAt, &pr.Pending)
		if err != nil {
			return nil, err
		}
		postReports = append(postReports, pr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return postReports, nil
}

func (mr *ModeratorDBRepository) ApprovePost(postID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE posts
						 SET is_approved = 1
						 WHERE id = ? 
		`, postID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (mr *ModeratorDBRepository) GetAllUnapprovedPosts() (posts []models.Post, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT *
							 FROM posts
							 WHERE is_approved = 0
		`); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt, &p.EditedAt, &p.IsImage,
			&p.ImagePath, &p.IsApproved)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (mr *ModeratorDBRepository) BanPost(postID int64, bans []string) (err error) {
	var (
		ctx     context.Context
		tx      *sql.Tx
		banRepo = postRepo.NewBanDBRepository(mr.dbConn)
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE posts
						 SET is_banned = 1
						 WHERE id = ? 
		`, postID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	if err = banRepo.Create(postID, bans); err != nil {
		return err
	}
	return nil
}

func (mr *ModeratorDBRepository) GetPostReportByID(postReportID int64) (*models.PostReport, error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		pr  models.PostReport
		err error
	)
	ctx = context.Background()
	if tx, err = mr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT *
						  FROM post_reports
						  WHERE id = ?
	`, postReportID).Scan(&pr.ID, &pr.ModeratorID, &pr.PostID, &pr.CreatedAt, &pr.Pending); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &pr, nil
}
