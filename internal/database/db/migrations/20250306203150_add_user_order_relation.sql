-- +goose Up
-- +goose StatementBegin
ALTER TABLE Orders ADD CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Orders DROP CONSTRAINT fk_orders_users;
-- +goose StatementEnd
