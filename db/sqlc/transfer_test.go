package db

import (
	"context"
	"testing"

	"github.com/amit16110/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	requires := require.New(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	requires.NoError(err)
	requires.NotEmpty(transfer)
	requires.Equal(arg.FromAccountID, transfer.FromAccountID)
	requires.Equal(arg.ToAccountID, transfer.ToAccountID)
	requires.Equal(arg.Amount, transfer.Amount)

	requires.NotZero(transfer.ID)
	requires.NotZero(transfer.CreatedAt)

	return transfer
}
