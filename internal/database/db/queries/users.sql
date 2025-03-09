-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByLogin :one
SELECT * FROM users
WHERE login = $1;

-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUserBalance :one
UPDATE users set balance=$2 where id=$1 RETURNING *;

-- name: UpdateUserRole :one
UPDATE users set role=$2 where id=$1 RETURNING *;