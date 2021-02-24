CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	username TEXT UNIQUE,
	password BLOB,
	email TEXT UNIQUE,
	created_at INTEGER,
	last_active INTEGER,
	session_id TEXT,
	expires_at INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS users_username ON users (username);

CREATE INDEX IF NOT EXISTS users_cover ON users (username, password, email, session_id);

CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);

CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	author_id INTEGER,
	title TEXT,
	content TEXT,
	created_at INTEGER,
	edited_at INTEGER,
	is_image INTEGER,
	image_path TEXT,
	FOREIGN KEY (author_id) REFERENCES users (id) ON
DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts_categories_bridge (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER,
	category_id INTEGER,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON
DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories (id) ON
DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_rating (
	id INTEGER,
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
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	author_id INTEGER,
	post_id INTEGER,
	content TEXT,
	created_at INTEGER,
	edited_at INTEGER,
	FOREIGN KEY (author_id) REFERENCES users (id) ON
DELETE CASCADE,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON
DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	receiver_id INTEGER,
	post_id INTEGER,
	rate_id INTEGER,
	comment_id INTEGER,
	comment_rate_id INTEGER,
	created_at INTEGER,
	FOREIGN KEY (receiver_id) REFERENCES users (id) ON
DELETE CASCADE,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON
DELETE CASCADE,
	FOREIGN KEY (comment_id) REFERENCES comments (id) ON
DELETE CASCADE,
	FOREIGN KEY (rate_id) REFERENCES post_rating (id) ON
DELETE CASCADE,
	FOREIGN KEY (comment_rate_id) REFERENCES comment_rating (id) ON
DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_rating (
	id INTEGER,
	user_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
	comment_id INTEGER NOT NULL,
	rate INTEGER,
	FOREIGN KEY (comment_id) REFERENCES comments (id) ON
DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users (id) ON
DELETE CASCADE,
	FOREIGN KEY(post_id) REFERENCES posts (id) ON
DELETE CASCADE,
	PRIMARY KEY(id)
);