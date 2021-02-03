package repository

import (
	"database/sql"
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
	"net/http"
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
func (pr *PostDBRepository) GetPostsByCategories(categories []string, userID int64) (posts []models.Post, status int, err error) {
	return nil, 0, nil
}

func (pr *PostDBRepository) GetPostsByRating(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	return nil, 0, nil
}

func (pr *PostDBRepository) GetPostsByDate(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	return nil, 0, nil
}
