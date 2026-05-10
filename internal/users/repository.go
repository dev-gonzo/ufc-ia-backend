package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"ufc-backend/internal/auth"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *User) error {
	query := `
		INSERT INTO users (
			email,
			username,
			password_hash,
			role
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return ErrEmailInUse
		case "users_username_key":
			return ErrUsernameInUse
		}
	}

	return err
}

func (r *Repository) FindByEmail(
	email string,
) (*auth.User, error) {

	query := `
		SELECT
			id,
			email,
			password_hash,
			role
		FROM users
		WHERE email = $1
	`

	var user auth.User

	err := r.db.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindInternalByEmail(
	email string,
) (*User, error) {

	query := `
		SELECT
			id,
			email,
			username,
			password_hash,
			role,
			is_active,
			created_at,
			last_login_at
		FROM users
		WHERE email = $1
	`

	var user User

	err := r.db.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.LastLoginAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) List() ([]User, error) {
	query := `
		SELECT
			id,
			email,
			username,
			role,
			is_active,
			created_at,
			last_login_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.LastLoginAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) UpdatePassword(
	userID string,
	passwordHash string,
) error {

	query := `
		UPDATE users
		SET
			password_hash = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		passwordHash,
		userID,
	)

	return err
}

func (r *Repository) UpdateRole(
	userID string,
	role string,
) error {

	query := `
		UPDATE users
		SET
			role = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		role,
		userID,
	)

	return err
}
