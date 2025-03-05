CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    number TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL,
    accrual FLOAT,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);
