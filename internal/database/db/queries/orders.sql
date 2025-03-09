-- name: GetOrderByID :one
SELECT * FROM Orders
WHERE id = $1 LIMIT 1;

-- name: GetOrdersByUser :many
SELECT * FROM Orders
WHERE user_id = $1
ORDER BY created_at;

-- name: CreateOrder :one
INSERT INTO Orders (id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateOrderStatus :exec
UPDATE Orders SET status=$3 WHERE id=$1 and user_id=$2;
