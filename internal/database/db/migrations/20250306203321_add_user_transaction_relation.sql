-- +goose Up
-- +goose StatementBegin
ALTER TABLE Transactions ADD CONSTRAINT fk_user_transactions FOREIGN KEY (user_id) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Transactions DROP CONSTRAINT fk_user_transactions;
-- +goose StatementEnd
