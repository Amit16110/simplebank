package db

import (
	"context"
	"testing"

	"github.com/amit16110/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	requires := require.New(t)
	requires.NoError(err)
	requires.NotEmpty(account)

	requires.Equal(arg.Owner, account.Owner)
	requires.Equal(arg.Balance, account.Balance)
	requires.Equal(arg.Currency, account.Currency)

	requires.NotZero(account.ID)
	requires.NotZero(account.CreatedAt)

}
