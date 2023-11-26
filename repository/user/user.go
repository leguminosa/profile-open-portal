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

// New returns a new instance of UserRepository.
func New(opts NewRepositoryOptions) *UserRepository {
	return &UserRepository{
		db: opts.DB,
	}
}

// GetUserByPhoneNumber returns a single user because phone number is stored unqiuely.
func (r *UserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	var user = &entity.User{}

	query := `
		SELECT
			id,
			fullname,
			phone_number,
			password,
			login_count,
			created_at,
			COALESCE(updated_at, created_at) AS updated_at
		FROM users
		WHERE phone_number = $1;
	`
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.Fullname,
		&user.PhoneNumber,
		&user.HashedPassword,
		&user.LoginCount,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID returns a single user by its id.
func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	var user = &entity.User{}

	query := `
		SELECT
			id,
			fullname,
			phone_number,
			password,
			login_count,
			created_at,
			COALESCE(updated_at, created_at) AS updated_at
		FROM users
		WHERE id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Fullname,
		&user.PhoneNumber,
		&user.HashedPassword,
		&user.LoginCount,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
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

// UpdateUser only updates fullname and phone number of a user with given id.
func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		UPDATE users
		SET
			fullname = $1,
			phone_number = $2
			updated_at = now()
		WHERE id = $3;
	`
	_, err = tx.ExecContext(
		ctx,
		query,
		user.Fullname,
		user.PhoneNumber,
		user.ID,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// IncrementLoginCount adds the value by 1 each time user logged in successfully.
func (r *UserRepository) IncrementLoginCount(ctx context.Context, userID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		UPDATE users
		SET
			login_count = login_count + 1,
			updated_at = now()
		WHERE id = $1;
	`
	_, err = r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
