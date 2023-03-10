// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) error
	CreateAccount(ctx context.Context, arg CreateAccountParams) (sql.Result, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (sql.Result, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (sql.Result, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountBalance(ctx context.Context, id int64) (int64, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetEntryForUpdate(ctx context.Context, id int64) (Entry, error)
	GetSession(ctx context.Context, id []byte) (Session, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	GetTransferForUpdate(ctx context.Context, id int64) (Transfer, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListAccount(ctx context.Context, arg ListAccountParams) ([]Account, error)
	ListEntry(ctx context.Context, arg ListEntryParams) ([]Entry, error)
	ListTransfer(ctx context.Context, arg ListTransferParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
