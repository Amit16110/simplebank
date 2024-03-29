package db

import (
	"context"
	"testing"

	"github.com/amit16110/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
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

	return account
}

func TestCreateAccount(t *testing.T) {

	createRandomAccount(t)
}
