package config

import (
	"time"
)

const (
	// URI
	APIPortDev    = "8081"
	APIURLDev     = "localhost:" + APIPortDev
	ClientURLDev  = "https://localhost:3000"
	ClientURLProd = "https://forume-react.herokuapp.com"

	// Session
	SessionCookieName = "forumSecretKey"
	SessionExpiration = 1 * time.Hour

	// User roles
	RoleGuest     = -1
	RoleUser      = 0
	RoleModerator = 1
	RoleAdmin     = 2

	// Database
	DBDriver   = "sqlite3"
	DBPath     = "./db"
	DBFileName = "forum.db"
	DBSchema   = "schema.sql"

	// Images
	ImagesPath   = "./images"
	MaxImageSize = 20 * 1024 * 1024
)

var (
	ClientURL = ClientURLDev
)
