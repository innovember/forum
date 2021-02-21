package config

import (
	"time"
)

const (
	APIPortDev        = "8081"
	APIURLDev         = "localhost:" + APIPortDev
	DBDriver          = "sqlite3"
	DBPath            = "./db"
	ImagesPath        = "./images"
	MaxImageSize      = 20 * 1024 * 1024
	DBFileName        = "forum.db"
	DBSchema          = "schema.sql"
	SessionCookieName = "forumSecretKey"
	SessionExpiration = 1 * time.Hour
	ClientURLDev      = "http://localhost:3000"
	ClientURLProd     = "https://forume-react.herokuapp.com"
)

var (
	ClientURL = ClientURLDev
)
