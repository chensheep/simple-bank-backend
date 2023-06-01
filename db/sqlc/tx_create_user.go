package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams CreateUserParams
	AfterUserCreated func(User) error
}

type CreateUserTxResult struct {
	User User `json:"user"`
}

func (s *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		newUser, err := q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		if err := arg.AfterUserCreated(newUser); err != nil {
			return err
		}

		result.User = newUser

		return nil
	})

	return result, err
}
