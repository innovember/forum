package config

import (
	"time"
)

const (
	APIPortDev        = "8081"
	APIURLDev         = "http://localhost:" + APIPortDev
	DBDriver          = "sqlite3"
	DBPath            = "./db"
	DBFileName        = "forum.db"
	DBSchema          = "schema.sql"
	SessionCookieName = "forumSecretKey"
	SessionExpiration = 1 * time.Hour
	ClientURL         = "http://localhost:3000"
)
