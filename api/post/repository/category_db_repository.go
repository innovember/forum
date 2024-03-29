package repository

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
)

type CategoryDBRepository struct {
	dbConn *sql.DB
}

func NewCategoryDBRepository(conn *sql.DB) post.CategoryRepository {
	return &CategoryDBRepository{dbConn: conn}
}

func (cr *CategoryDBRepository) Create(postID int64, categories []string) (err error) {
	var (
		categoryID int64
		result     sql.Result
		isExist    bool
	)
	for _, category := range categories {
		if isExist, err = cr.IsCategoryExist(category); err != nil {
			return err
		}
		if !isExist {
			if result, err = cr.dbConn.Exec(`INSERT INTO categories(name) VALUES(?)`, category); err != nil {
				return err
			}
			if categoryID, err = result.LastInsertId(); err != nil {
				return err
			}
		} else {
			if categoryID, err = cr.GetCategoryIDByName(category); err != nil {
				return err
			}
		}
		if _, err = cr.dbConn.Exec(
			`INSERT INTO posts_categories_bridge (post_id, category_id)
			VALUES (?, ?)`,
			postID, categoryID,
		); err != nil {
			return err
		}
	}
	return nil
}

func (cr *CategoryDBRepository) IsCategoryExist(category string) (bool, error) {
	var (
		id  int64
		err error
	)
	if err = cr.dbConn.QueryRow(`SELECT id FROM categories WHERE name=?`, category).Scan(
		&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (cr *CategoryDBRepository) GetCategoryIDByName(name string) (id int64, err error) {
	if err = cr.dbConn.QueryRow(`SELECT id FROM categories WHERE name=?`, name).Scan(
		&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return id, nil
}

func (cr *CategoryDBRepository) GetAllCategories() (categories []models.Category, status int, err error) {
	var rows *sql.Rows
	if rows, err = cr.dbConn.Query(`SELECT * FROM categories`); err != nil {
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

func (cr *CategoryDBRepository) Update(postID int64, categories []string) (err error) {
	if err = cr.DeleteFromPostCategoriesBridge(postID); err != nil {
		return err
	}
	if err = cr.Create(postID, categories); err != nil {
		return err
	}
	return nil
}

func (cr *CategoryDBRepository) DeleteFromPostCategoriesBridge(postID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM posts_categories_bridge
						WHERE post_id = ?`, postID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM categories
						 WHERE id IN
						(SELECT c.id FROM categories AS c 
						LEFT JOIN posts_categories_bridge AS pcb
						ON c.id = pcb.category_id
						WHERE pcb.category_id IS NULL
 						)`); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *CategoryDBRepository) DeleteCategoryByID(categoryID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM categories
						 WHERE id = ?
 						`, categoryID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *CategoryDBRepository) CreateNewCategory(category string) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO categories(name) VALUES(?)`, category); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
