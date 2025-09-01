package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandAccount(t)
	account2 := createRandAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	// run 5 concurrent transactions
	n := 5
	amount := int64(10)

	// two channels for getting the result of go routines
	errs := make(chan error)
	results := make(chan TransferTxResults)

	// do 5 transactions as go routines pass in the return err and result over channel to the main
	for range n{
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		} ()
	}

	existed := make(map[int]bool)
	// for each transaction, check the results and errors
	for range n{

		// check for no errors
		err := <- errs
		require.NoError(t, err)

		result := <- results

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, errTx := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, errTx)
		
		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)

		//check balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1>0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// get account and check balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	// check balance
	require.Equal(t, account1.Balance - int64(n) * amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandAccount(t)
	account2 := createRandAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	// run 10 concurrent transactions (odd are from 1 to 2 even from 2 to 1)
	n := 10
	amount := int64(10)

	// two channels for getting the result of go routines
	errs := make(chan error)
	// results := make(chan TransferTxResults) checked in previous tests

	// do 5 transactions as go routines pass in the return err and result over channel to the main
	for i := 0; i < n; i++{
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 == 1{
			fromAccountID = account2.ID
			toAccountID  = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})
			errs <- err
		} ()
	}

	// for each transaction, check the results and errors
	for i:= 0; i < n; i++{
		// check for no errors
		err := <- errs
		require.NoError(t, err)
	}

	// get account and check balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	// check balance
	require.Equal(t, account1.Balance, updatedAccount1.Balance) // same balance as acc1 - 2 5tx ans 2 - 1 5tx
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}