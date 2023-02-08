-- name: CreateEntry :execresult
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  ?, ?
);

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ? LIMIT 1;

-- name: GetEntryForUpdate :one
SELECT * FROM entries
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: ListEntry :many
SELECT * FROM entries
WHERE account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?;