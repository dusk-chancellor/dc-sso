package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/dusk-chancellor/dc-sso/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// db query calls implementation

// create new user, returning `id`
func (r *Repo) CreateUser(ctx context.Context, username, email string, password []byte) (string, error) {
	q := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	var id string

	if err := r.db.QueryRow(ctx, q, username, email, password).Scan(&id); err != nil {
		var e *pgconn.PgError
		// if fields with unique property repeat
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return "", ErrUserAlreadyExists
		}
		// other errors
		return "", err
	}

	return id, nil
}

// retrieve user by {id / email / username}, returning user model
func (r *Repo) GetUserBy(ctx context.Context, field Field, val any) (models.User, error) {
	q := fmt.Sprintf(`
	SELECT id, username, email, password, role
	FROM users
	WHERE %s = $1
	`, field)

	var user models.User

	if err := r.db.QueryRow(ctx, q, val).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		// if no rows w/ such {field} found
		if err == pgx.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		// other errors
		return models.User{}, err
	}

	return user, nil
}
// update user {email, username, password} by `id`, returning success
func (r *Repo) UpdateByID(ctx context.Context, id string, field Field, val any) error {
	q := fmt.Sprintf(`
	UPDATE users
	SET %s = $2
	WHERE id = $1;
	`, field)

	if _, err := r.db.Exec(ctx, q, id, val); err != nil {
		var e *pgconn.PgError
		// if fields with unique property repeat
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return ErrUserAlreadyExists
		}
		// if no rows w/ such `id` found
		if err == pgx.ErrNoRows {
			return ErrUserNotFound
		}
		// other errors
		return err
	}

	return nil
}

func (r *Repo) UpdateRole(ctx context.Context, email, wantsRole string) error {
	q := `
	UPDATE users
	SET role = $2
	WHERE email = $1;
	`

	if _, err := r.db.Exec(ctx, q, email, wantsRole); err != nil {
		// if no row w/ such `email` found
		if err == pgx.ErrNoRows {
			return ErrUserNotFound
		}
		// other errors
		return err
	}

	return nil
}

