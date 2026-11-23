package store

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id 			  INTEGER PRIMARY KEY,
	email 		  TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	is_admin 	  INTEGER NOT NULL DEFAULT 0,
	confirmed 	  INTEGER NOT NULL DEFAULT 0,
	confirm_code  TEXT,
	reset_code 	  TEXT,
	created_at 	  TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
	token	   TEXT PRIMARY KEY,
	user_id	   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
