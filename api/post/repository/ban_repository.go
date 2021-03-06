package repository

import (
	"context"
	"database/sql"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
)

type BanDBRepository struct {
	dbConn *sql.DB
}

func NewBanDBRepository(conn *sql.DB) post.BanRepository {
	return &BanDBRepository{dbConn: conn}
}

func (br *BanDBRepository) Create(postID int64, categories []string) (err error) {
	var (
		categoryID int64
		result     sql.Result
		isExist    bool
	)
	for _, category := range categories {
		if isExist, err = br.IsCategoryExist(category); err != nil {
			return err
		}
		if !isExist {
			if result, err = br.dbConn.Exec(`INSERT INTO bans(name) VALUES(?)`, category); err != nil {
				return err
			}
			if categoryID, err = result.LastInsertId(); err != nil {
				return err
			}
		} else {
			if categoryID, err = br.GetCategoryIDByName(category); err != nil {
				return err
			}
		}
		if _, err = br.dbConn.Exec(
			`INSERT INTO posts_bans_bridge (post_id, ban_id)
			VALUES (?, ?)`,
			postID, categoryID,
		); err != nil {
			return err
		}
	}
	return nil
}

func (br *BanDBRepository) IsCategoryExist(category string) (bool, error) {
	var (
		id  int64
		err error
	)
	if err = br.dbConn.QueryRow(`SELECT id FROM bans WHERE name=?`, category).Scan(
		&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (br *BanDBRepository) GetCategoryIDByName(name string) (id int64, err error) {
	if err = br.dbConn.QueryRow(`SELECT id FROM bans WHERE name=?`, name).Scan(
		&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return id, nil
}

func (br *BanDBRepository) GetAllCategories() (categories []models.Category, status int, err error) {
	var rows *sql.Rows
	if rows, err = br.dbConn.Query(`SELECT * FROM bans`); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Category
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		categories = append(categories, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return categories, http.StatusOK, nil
}

func (br *BanDBRepository) Update(postID int64, categories []string) (err error) {
	if err = br.DeleteFromPostCategoriesBridge(postID); err != nil {
		return err
	}
	if err = br.Create(postID, categories); err != nil {
		return err
	}
	return nil
}

func (br *BanDBRepository) DeleteFromPostCategoriesBridge(postID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = br.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM posts_bans_bridge
						WHERE post_id = ?`, postID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM bans
						 WHERE id IN
						(SELECT b.id FROM bans AS b 
						LEFT JOIN posts_bans_bridge AS pbb
						ON b.id = pbb.ban_id
						WHERE pbb.ban_id IS NULL
 						)`); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (br *BanDBRepository) DeleteCategoryByID(categoryID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = br.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM bans
						 WHERE id = ?
 						)`, categoryID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (br *BanDBRepository) CreateNewCategory(category string) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = br.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO bans(name) VALUES(?)`, category); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
