package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/chensheep/simple-bank-backend/db/mock"
	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/chensheep/simple-bank-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {

	account := createRandomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "StatusBadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStore := mockdb.NewMockStore(mockCtrl)
			// build stubs
			tc.buildStubs(mockStore)

			server := NewServer(mockStore)
			w := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", tc.accountID), nil)
			require.NoError(t, err)

			server.router.ServeHTTP(w, r)
			// check response
			tc.checkResponse(t, w)
		})
	}
}

func TestCreateAccount(t *testing.T) {

	account := createRandomAccount()

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		body          gin.H
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0,
					}).Times(1).
					Return(account, nil)
			},
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalServerError",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			body: gin.H{
				"owner":    account.Owner,
				"currency": "invalid",
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStore := mockdb.NewMockStore(mockCtrl)
			// build stubs
			tc.buildStubs(mockStore)

			server := NewServer(mockStore)
			w := httptest.NewRecorder()

			arg, err := json.Marshal(tc.body)
			require.NoError(t, err)

			httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(arg))
			server.router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(arg)))
			//check response
			tc.checkResponse(t, w)
		})
	}
}

// TODO: add test for list accounts

func createRandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	account2 := db.Account{}
	err = json.Unmarshal(data, &account2)
	require.NoError(t, err)
	require.Equal(t, account, account2)

	fmt.Println(account)
	fmt.Println(account2)
}
