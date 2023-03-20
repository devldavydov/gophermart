package storage

const (
	// users
	_sqlCreateTableUser = `
	CREATE TABLE IF NOT EXISTS users (
		id         serial                   NOT NULL,
		username   text                     NOT NULL,
		password   text                     NOT NULL,
		created_at timestamp with time zone NOT NULL DEFAULT now(),
		
		PRIMARY KEY (id),
		UNIQUE(username)
	);
	`
	_sqlCreateUser = `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id
	`
	_sqlFindUser = `
	SELECT id, password FROM users
	WHERE username = $1
	`
	// orders
	_sqlCreateTableOrders = `
	CREATE TABLE IF NOT EXISTS orders (
        number      text                     NOT NULL,
		user_id     int                      NOT NULL,
		status      text                     NOT NULL DEFAULT 'NEW',
		accrual     int,
		uploaded_at timestamp with time zone NOT NULL DEFAULT now(),
	
		PRIMARY KEY (number),
		FOREIGN KEY(user_id) REFERENCES users(id),
		CHECK(status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED'))
	);
	`
	_sqlAddOrder = `
	INSERT INTO orders (number, user_id)
	VALUES ($1, $2)
	`
	_sqlFindUserOrder = `
	SELECT 1 FROM orders
	WHERE number = $1 AND user_id = $2
	`
	_sqlListOrders = `
	SELECT number, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at ASC
	`
	// balance
	_sqlCreateTableBalance = `
	CREATE TABLE IF NOT EXISTS balance (
		user_id   int              NOT NULL,
		current   double precision NOT NULL DEFAULT 0,
		withdrawn double precision NOT NULL DEFAULT 0,

		PRIMARY KEY (user_id),
		FOREIGN KEY(user_id) REFERENCES users(id),
		CHECK(current >= 0)
	);
	`
	_sqlCreateUserBalance = `
	INSERT INTO balance (user_id)
	VALUES ($1)
	`
	_sqlGetUserBalance = `
	SELECT current, withdrawn FROM balance
	WHERE user_id = $1
	`
)
