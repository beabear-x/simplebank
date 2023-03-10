package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/beabear/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) int64 {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	result, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	accountId, err := result.LastInsertId()
	require.NoError(t, err)
	return accountId
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountId1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), accountId1)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
}

func TestUpdateAccount(t *testing.T) {
	accountId1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		Balance: util.RandomMoney(),
		ID:      accountId1,
	}

	err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	accountId1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), accountId1)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), accountId1)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	var lastAccountID int64
	for i := 0; i < 10; i++ {
		lastAccountID = createRandomAccount(t)
	}
	lastAccount, err := testQueries.GetAccount(context.Background(), lastAccountID)
	require.NoError(t, err)

	arg := ListAccountParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
