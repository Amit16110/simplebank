package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries //Note:- Composition => it work like a inheritance in golang.
	db       *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database	transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return fmt.Errorf("Tx err: %v, rollback err: %v", err, rollBackErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id,omitempty"`
	ToAccountID   int64 `json:"to_account_id,omitempty"`
	Amount        int64 `json:"amount,omitempty"`
}

//TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer,omitempty"`
	FromAccount Account  `json:"from_account,omitempty"`
	ToAccount   Account  `json:"to_account,omitempty"`
	FromEntry   Entry    `json:"from_entry,omitempty"`
	ToEntry     Entry    `json:"to_entry,omitempty"`
}

// Transfer performs a money trasfer from one account to the other.
// It creates a transfer record, add account entries, and updates accounts balance within a single database transaction

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// if arg.FromAccountID < arg.ToAccountID {
		// 	result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		// } else {
		// 	result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		// }

		return nil

	})
	return result, err
}

// func addMoney(ctx context.Context, q *Queries, accountID1, ammount1, accountID2, amount2 int64) (account1, account2, Account, err error) {
// 	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
// 		ID:     accountID1,
// 		Amount: amount1,
// 	})
// 	if err != nil {
// 		return
// 	}

// 	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
// 		ID:     accountID2,
// 		Amount: amount2,
// 	})
// 	return
// }
