-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Transactions(
    id serial PRIMARY KEY,
    user_id integer not null,
    order_id integer not null,
    type_transaction VARCHAR(15) not null default 'ACCRUAL' CHECK (type_transaction IN ('WITHDRAWAL', 'ACCRUAL')),
    amount float not null default 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Transactions;
-- +goose StatementEnd
