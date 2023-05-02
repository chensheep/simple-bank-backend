package db

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	errChan := make(chan error)
	resChan := make(chan TransferTxResult)
	var wg sync.WaitGroup

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(errChan chan error, resChan chan TransferTxResult) {
			defer wg.Done()
			result, err := store.TransferTx(context.Background(), arg)
			errChan <- err
			resChan <- result

		}(errChan, resChan)
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err, res := <-errChan, <-resChan

		require.NoError(t, err)
		require.NotEmpty(t, res)

		require.NotEmpty(t, res.Transfer)
		require.Equal(t, arg.FromAccountID, res.Transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, res.Transfer.ToAccountID)
		require.Equal(t, arg.Amount, res.Transfer.Amount)
		require.NotZero(t, res.Transfer.ID)
		require.NotZero(t, res.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		require.NotEmpty(t, res.FromEntry)
		require.Equal(t, arg.FromAccountID, res.FromEntry.AccountID)
		require.Equal(t, -arg.Amount, res.FromEntry.Amount)
		require.NotZero(t, res.FromEntry.ID)
		require.NotZero(t, res.FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), res.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, res.ToEntry)
		require.Equal(t, arg.ToAccountID, res.ToEntry.AccountID)
		require.Equal(t, arg.Amount, res.ToEntry.Amount)
		require.NotZero(t, res.ToEntry.ID)
		require.NotZero(t, res.ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), res.ToEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%arg.Amount == 0)

		k := int(diff1 / arg.Amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, updateAccount1.Balance, account1.Balance-int64(n)*arg.Amount)

	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updateAccount2.Balance, account2.Balance+int64(n)*arg.Amount)

	wg.Wait()

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	errChan := make(chan error)

	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		var arg TransferTxParams
		if i%2 == 0 {
			arg = TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        10,
			}
		} else {
			arg = TransferTxParams{
				FromAccountID: account2.ID,
				ToAccountID:   account1.ID,
				Amount:        10,
			}
		}
		go func(errChan chan error, arg TransferTxParams) {
			defer wg.Done()

			_, err := store.TransferTx(context.Background(), arg)
			errChan <- err
		}(errChan, arg)
	}

	for i := 0; i < n; i++ {
		err := <-errChan

		require.NoError(t, err)
	}

	// check the final updated balance
	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, updateAccount1.Balance, account1.Balance)

	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updateAccount2.Balance, account2.Balance)

	wg.Wait()

}
