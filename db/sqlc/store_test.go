package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	//NewStore require a db object so
	//Make a Db object global

	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("data")
	t.Log("show data", account1.Balance)
	t.Log("account>>", account1.Balance, account2.Balance)

	//  In database transaction, we need to handle concurrecy carefully.

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	erros := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			erros <- err
			results <- result
		}()
	}

	// Check result
	for i := 0; i < n; i++ {
		err := <-erros
		if err != nil {
			t.Error(err)
		}
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		fmt.Println("resuilt", result.Transfer)
		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntries := result.FromEntry
		require.NotEmpty(t, fromEntries)
		require.Equal(t, account1.ID, fromEntries.AccountID)
		require.Equal(t, -amount, fromEntries.Amount)
		require.NotZero(t, fromEntries.ID)
		require.NotZero(t, fromEntries.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntries.ID)
		require.NoError(t, err)

		toEntries := result.ToEntry
		require.NotEmpty(t, toEntries)
		require.Equal(t, account2.ID, toEntries.AccountID)
		require.Equal(t, amount, toEntries.Amount)
		require.NotZero(t, toEntries.ID)
		require.NotZero(t, toEntries.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntries.ID)
		require.NoError(t, err)

		// Todo: check the account's balance
	}
	fmt.Println("account>>", account1.Balance, account2.Balance)
}
