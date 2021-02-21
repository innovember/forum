package db

import (
	"database/sql"
	"fmt"
	config "github.com/innovember/forum/api/config"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"strings"
)

var (
	DBConn *sql.DB
	err    error
)

// get instance of db connection, and check db integrity with schema
func GetDBInstance() (*sql.DB, error) {
	var (
		DB_USER string = os.Getenv("DB_USER")
		DB_PASS string = os.Getenv("DB_PASS")
		DB_AUTH string
		DB_URI  string
	)
	if _, err = os.Stat(config.DBPath); os.IsNotExist(err) {
		if err = os.Mkdir(config.DBPath, 0755); err != nil {
			return nil, err
		}
	}
	if _, err = os.Stat(config.ImagesPath); os.IsNotExist(err) {
		if err = os.Mkdir(config.ImagesPath, 0755); err != nil {
			return nil, err
		}
	}
	if DB_PASS != "" && DB_USER != "" {
		DB_AUTH = fmt.Sprintf("?_auth&_auth_user=%s&_auth_pass=%s", DB_USER, DB_PASS)
	}
	DB_URI = fmt.Sprintf("%s/%s%s", config.DBPath, config.DBFileName, DB_AUTH)
	if DBConn, err = sql.Open(config.DBDriver, DB_URI); err != nil {
		return nil, err
	}
	DBConn.SetMaxIdleConns(100)
	if err = DBConn.Ping(); err != nil {
		return nil, err
	}
	return DBConn, nil
}

func CheckDB(DBConn *sql.DB, path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	queries := strings.Split(string(schema), ";\n")
	for _, query := range queries {
		_, err = DBConn.Exec(string(query))
		if err != nil {
			return err
		}
	}
	return nil
}
