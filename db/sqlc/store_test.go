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
	fmt.Println("data===>", account1.Balance, account2.Balance)

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
	existed := make(map[int]bool)
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

		// check accounts
		fromAccount := result.FromAccount
		fmt.Println("fromAccount", fromAccount)
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Todo: check the account's balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	fmt.Println("account>>", account1.Balance, account2.Balance)
	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
