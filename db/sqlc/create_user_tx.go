package db

import "context"

type CreateUserTxParams struct {
	CreateuserParams
	AterTx func(user User) error
}

type CreateUserTxResult struct {
	User
}

func (store *Store) CreateUserTx(ctx context.Context, args CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.Createuser(ctx, args.CreateuserParams)
		if err != nil {
			return err
		}

		return args.AterTx(result.User)
	})

	return result, err
}
