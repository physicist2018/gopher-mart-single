-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Orders(
    id INTEGER PRIMARY KEY,
    user_id integer not null,
    status VARCHAR(20) not null default 'NEW' CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED', 'REGISTERED')),
    accrual float not null default 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Orders;
-- +goose StatementEnd
