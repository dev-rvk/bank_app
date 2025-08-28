package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store strucct
type Store struct{
	*Queries
	db *sql.DB
}

// constructor for a new store
func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// base function to perform transaction given the queries
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}


type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResults struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`

}

// function used to transfer money between two accounts
// 1) Create transfer record, 2) Add account entries, 3) update account balance
func (store *Store) TransferTx (ctx context.Context, args TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err :=  store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID: args.ToAccountID,
			Amount: args.Amount,
		})

		if err != nil {
			return err
		}

		// create entry in both accounts
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount: -args.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: update account balance to prevent deadlock

		return nil
	})

	return result, err
}