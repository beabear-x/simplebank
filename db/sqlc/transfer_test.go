package db

import (
	"context"
	"testing"

	"github.com/beabear/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, accountId1, accountId2 int64) int64 {

	arg := CreateTransferParams{
		FromAccountID: accountId1,
		ToAccountID:   accountId2,
		Amount:        util.RandomMoney(),
	}

	result, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	transferId, err := result.LastInsertId()
	require.NoError(t, err)
	require.NotEmpty(t, transferId)

	return transferId
}

func TestCreateTransfer(t *testing.T) {
	accountId1 := createRandomAccount(t)
	accountId2 := createRandomAccount(t)
	createRandomTransfer(t, accountId1, accountId2)
}

func TestGetTransfer(t *testing.T) {
	accountId1 := createRandomAccount(t)
	accountId2 := createRandomAccount(t)
	transferId := createRandomTransfer(t, accountId1, accountId2)

	transfer, err := testQueries.GetTransfer(context.Background(), transferId)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, accountId1, transfer.FromAccountID)
	require.Equal(t, accountId2, transfer.ToAccountID)
}

func TestListTransfer(t *testing.T) {
	accountId1 := createRandomAccount(t)
	accountId2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, accountId1, accountId2)
	}

	arg := ListTransferParams{
		FromAccountID: accountId1,
		ToAccountID:   accountId2,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfer(context.Background(), arg)

	require.NoError(t, err)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, accountId1, transfer.FromAccountID)
		require.Equal(t, accountId2, transfer.ToAccountID)
	}
}
