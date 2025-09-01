package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/devrvk/simplebank/db/mock"
	db "github.com/devrvk/simplebank/db/sqlc"
	"github.com/devrvk/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {

	// controller -> store -> server
	// create a request and serve it on the server

	// generate a random account
	account := randomAccount()

	// create a new controller (to create store)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	
}

// function to return a random account
func randomAccount () (db.Account){
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: util.RandOwner(),
		Balance: util.RandBalance(),
		Currency: util.RandCurrency(),
	}
}