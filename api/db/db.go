package db

import (
	"database/sql"
	config "github.com/innovember/forum/api/config"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"strings"
)

var DBConn *sql.DB
var err error

// get instance of db connection, and check db integrity with schema
func GetDBInstance() (*sql.DB, error) {
	if _, err = os.Stat(config.DBPath); os.IsNotExist(err) {
		if err = os.Mkdir(config.DBPath, 0755); err != nil {
			return nil, err
		}
	}
	if DBConn, err = sql.Open(config.DBDriver, config.DBPath+"/"+config.DBFileName); err != nil {
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
