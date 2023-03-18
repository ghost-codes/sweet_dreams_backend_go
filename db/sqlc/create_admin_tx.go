package db

import "context"

type CreateAdminTxParams struct {
	CreateAdminParams
	AfterCreate func(admin Admin) error
}

type CreateAdminTxResult struct {
	Admin Admin
}

func (store *Store) CreateAdminTX(ctx context.Context, args CreateAdminTxParams) (CreateAdminTxResult, error) {
	var result CreateAdminTxResult

	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		admin, err := store.CreateAdmin(ctx, args.CreateAdminParams)
		if err != nil {
			return err
		}
		result = CreateAdminTxResult{Admin: admin}
		return args.AfterCreate(admin)
	})

	return result, err
}
