-- +goose Up
-- +goose StatementBegin
ALTER TABLE Transactions ADD CONSTRAINT unq_transaction_unique_columns UNIQUE (user_id, order_id, type_transaction);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Transactions DROP CONSTRAINT unq_transaction_unique_columns
-- +goose StatementEnd
