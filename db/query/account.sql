-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  ?, ?, ?
);

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: GetAccountBalance :one
SELECT balance FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :exec
UPDATE accounts 
SET balance = ?
WHERE id = ?;

-- name: AddAccountBalance :exec
UPDATE accounts 
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;