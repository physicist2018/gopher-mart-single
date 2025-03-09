-- +goose Up
-- +goose StatementBegin
ALTER TABLE Transactions ADD CONSTRAINT fk_orders_transactions FOREIGN KEY (order_id) REFERENCES Orders(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Transactions DROP CONSTRAINT fk_orders_transactions;
-- +goose StatementEnd
