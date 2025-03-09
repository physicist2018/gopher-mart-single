-- name: CreateTransaction :one
INSERT INTO Transactions (user_id, order_id, type_transaction, amount) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateWithdrawalTransaction :exec
INSERT INTO Transactions(user_id, order_id, type_transaction, amount)
VALUES ($1, $2, 'WITHDRAWAL', $3);

-- name: CreateAccrualTransaction :exec
INSERT INTO Transactions(user_id, order_id, type_transaction, amount)
VALUES ($1, $2, 'ACCRUAL', $3);
