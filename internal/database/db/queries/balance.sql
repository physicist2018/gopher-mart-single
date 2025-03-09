-- name: GetUserBalance :one
SELECT 
    u.balance + 
    COALESCE(SUM(CASE WHEN t.type_transaction = 'ACCRUAL' THEN t.amount ELSE 0 END), 0) - 
    COALESCE(SUM(CASE WHEN t.type_transaction = 'WITHDRAWAL' THEN t.amount ELSE 0 END), 0) AS current,
    CAST(COALESCE(SUM(CASE WHEN t.type_transaction = 'WITHDRAWAL' THEN t.amount ELSE 0 END), 0) AS FLOAT) AS withdrawn
FROM users u
LEFT JOIN Transactions t ON u.id = t.user_id
WHERE u.id = $1
GROUP BY u.id;