package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/devrvk/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandAccount(t *testing.T) Account {

	// testing data
	arg := CreateAccountParams{
		Owner: util.RandOwner(),
		Balance: util.RandBalance(),
		Currency: util.RandCurrency(),
	}

	// testQueries is struct obtained in main_testing.go on connection with the database
	// CreateAccount() is a function used in the testQueries struct which is executed for testing
	account, err := testQueries.CreateAccount(context.Background(), arg)

	
	// we use require to check for errors, compare values, itc from the testify package

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T){
	createRandAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {

	account1 := createRandAccount(t)

	args := UpdateAccountParams{
		ID: account1.ID,
		Balance: util.RandBalance(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Balance, args.Balance)
	require.Equal(t, account2.Currency, account1.Currency)
	require.Equal(t, account2.Owner, account1.Owner)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {
	for i:= 0; i < 10; i++ {
		createRandAccount(t)
	}

	arg := ListAccountParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts{
		require.NotEmpty(t, account)
	}

}