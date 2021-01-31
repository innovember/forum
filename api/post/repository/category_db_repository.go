package repository

import (
	"database/sql"
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
			return true, nil
		}
		return false, err
	}
	return false, nil
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
