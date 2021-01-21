CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			username TEXT UNIQUE,
			password BLOB,
			email TEXT UNIQUE,
			created_at INTEGER,
			last_active INTEGER,
			session_id TEXT
		);
CREATE INDEX IF NOT EXISTS users_username ON users (username);
CREATE INDEX IF NOT EXISTS users_cover ON users (username, password,email,session_id);
