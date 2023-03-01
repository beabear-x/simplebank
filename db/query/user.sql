-- name: CreateUser :execresult
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  ?, ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: UpdateUser :execresult
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at)
WHERE username = sqlc.arg(username);