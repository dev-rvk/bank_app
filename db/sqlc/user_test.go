package db

import (
	"context"
	"testing"
	"time"

	"github.com/devrvk/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandUser(t *testing.T) User {

	// testing data
	arg := CreateUserParams{
		Username: util.RandOwner(),
		HashedPassword: "secret",
		FullName: util.RandOwner(),
		Email: util.RandEmail(),
	}

	// testQueries is struct obtained in main_testing.go on connection with the database
	// CreateAccount() is a function used in the testQueries struct which is executed for testing
	user, err := testQueries.CreateUser(context.Background(), arg)

	
	// we use require to check for errors, compare values, itc from the testify package

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser (t *testing.T){
	createRandUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}