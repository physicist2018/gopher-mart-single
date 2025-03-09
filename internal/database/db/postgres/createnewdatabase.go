package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(10) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    balance float NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Orders(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED', 'REGISTERED')),
    accrual float NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Transactions(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    type_transaction VARCHAR(15) NOT NULL DEFAULT 'ACCRUAL' CHECK (type_transaction IN ('WITHDRAWAL', 'ACCRUAL')),
    amount float NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`

const constraints = `
ALTER TABLE Orders ADD CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE Transactions ADD CONSTRAINT fk_orders_transactions FOREIGN KEY (order_id) REFERENCES Orders(id);
ALTER TABLE Transactions ADD CONSTRAINT fk_user_transactions FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE Transactions ADD CONSTRAINT unq_transaction_unique_columns UNIQUE (user_id, order_id, type_transaction);
`

func NewDB(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	db.MustExec(schema)
	_, err = db.Exec(constraints)
	return db
}
