package storage

const (
	_sqlCreateTableUser = `
	CREATE TABLE IF NOT EXISTS users (
		id       serial    NOT NULL,
		username text      NOT NULL,
		password text      NOT NULL,
		created  timestamp NOT NULL DEFAULT now(),
		
		PRIMARY KEY (id),
		UNIQUE(username)
	);
	`
)
