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

CREATE INDEX IF NOT EXISTS users_cover ON users (username, password, email, session_id);

CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT);

CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY,
	author_id INTEGER title TEXT,
	content TEXT,
	created_at INTEGER,
	FOREIGN KEY (author_id) REFERENCES users (id) ON
DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts_categories_bridge (
	id INTEGER PRIMARY KEY,
	post_id INTEGER,
	category_id INTEGER,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON
DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories (id) ON
DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_rating (
	id INTEGER NOT NULL UNIQUE,
	user_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
	rate INTEGER,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON
DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users (id) ON
DELETE CASCADE,
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
	author_id INTEGER,
	post_id INTEGER,
	content TEXT,
	created_at INTEGER,
	FOREIGN KEY (author_id) REFERENCES posts (id) ON
DELETE CASCADE,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON
DELETE CASCADE
);