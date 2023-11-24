// This file contains the repository implementation layer.
package user

import (
	"context"
	"database/sql"

	"github.com/leguminosa/profile-open-portal/entity"
)

type UserRepository struct {
	db *sql.DB
}

type NewRepositoryOptions struct {
	DB *sql.DB
}

func New(opts NewRepositoryOptions) *UserRepository {
	return &UserRepository{
		db: opts.DB,
	}
}

// InsertUser inserts a new user to database, returning its id on success.
func (r *UserRepository) InsertUser(ctx context.Context, user *entity.User) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO users (
			fullname,
			phone_number,
			password
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING id;
	`
	err = tx.QueryRowContext(
		ctx,
		query,
		user.Fullname,
		user.PhoneNumber,
		user.HashedPassword,
	).Scan(&user.ID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
