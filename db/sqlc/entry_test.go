package db

import (
	"context"
	"testing"

	"github.com/beabear/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, accountId int64) int64 {
	arg := CreateEntryParams{
		AccountID: accountId,
		Amount:    util.RandomMoney(),
	}

	result, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	entryId, err := result.LastInsertId()
	require.NoError(t, err)
	require.NotEmpty(t, entryId)

	return entryId
}

func TestCreateEntry(t *testing.T) {
	accountId := createRandomAccount(t)
	createRandomEntry(t, accountId)
}

func TestGetEntry(t *testing.T) {
	accountId := createRandomAccount(t)
	entryId := createRandomEntry(t, accountId)
	entry, err := testQueries.GetEntry(context.Background(), entryId)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
}

func TestListEntry(t *testing.T) {
	accountId := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, accountId)
	}

	arg := ListEntryParams{
		AccountID: accountId,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
