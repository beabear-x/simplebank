// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: account.sql

package db

import (
	"context"
	"database/sql"
)

const addAccountBalance = `-- name: AddAccountBalance :exec
UPDATE accounts 
SET balance = balance + ?
WHERE id = ?
`

type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) error {
	_, err := q.db.ExecContext(ctx, addAccountBalance, arg.Amount, arg.ID)
	return err
}

const createAccount = `-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  ?, ?, ?
)
`

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = ? LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountBalance = `-- name: GetAccountBalance :one
SELECT balance FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE
`

func (q *Queries) GetAccountBalance(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getAccountBalance, id)
	var balance int64
	err := row.Scan(&balance)
	return balance, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listAccount = `-- name: ListAccount :many
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE owner = ?
ORDER BY id
LIMIT ?
OFFSET ?
`

type ListAccountParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListAccount(ctx context.Context, arg ListAccountParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccount, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :exec
UPDATE accounts 
SET balance = ?
WHERE id = ?
`

type UpdateAccountParams struct {
	Balance int64 `json:"balance"`
	ID      int64 `json:"id"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) error {
	_, err := q.db.ExecContext(ctx, updateAccount, arg.Balance, arg.ID)
	return err
}
