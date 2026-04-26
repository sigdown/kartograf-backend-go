package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

var _ usecase.UserRepository = (*PostgresUserRepository)(nil)

func (r *PostgresUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	row := r.pool.QueryRow(ctx, `
		insert into "user" (username, display_name, email, password_hash, role)
		values ($1, nullif($2, ''), $3, $4, $5)
		returning user_id, username, coalesce(display_name, ''), email, password_hash, role, created_at, updated_at
	`, user.Username, user.DisplayName, strings.ToLower(user.Email), user.PasswordHash, user.Role)

	created, err := scanUser(row)
	if err != nil {
		return domain.User{}, mapError(err, "user")
	}

	return created, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, userID int64) (domain.User, error) {
	row := r.pool.QueryRow(ctx, `
		select user_id, username, coalesce(display_name, ''), email, password_hash, role, created_at, updated_at
		from "user"
		where user_id = $1
	`, userID)

	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("%w: user not found", domain.ErrNotFound)
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByLogin(ctx context.Context, login string) (domain.User, error) {
	row := r.pool.QueryRow(ctx, `
		select user_id, username, coalesce(display_name, ''), email, password_hash, role, created_at, updated_at
		from "user"
		where lower(username) = $1 or lower(email) = $1
	`, strings.ToLower(login))

	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("%w: user not found", domain.ErrUnauthorized)
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, userID int64, input usecase.UpdateAccountInput) (domain.User, error) {
	setClauses := make([]string, 0, 4)
	args := []any{userID}
	idx := 2

	if input.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", idx))
		args = append(args, input.Username)
		idx++
	}

	if input.DisplayName != "" {
		setClauses = append(setClauses, fmt.Sprintf("display_name = $%d", idx))
		args = append(args, input.DisplayName)
		idx++
	}

	if input.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", idx))
		args = append(args, strings.ToLower(input.Email))
		idx++
	}

	if input.Password != "" {
		setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", idx))
		args = append(args, input.Password)
		idx++
	}

	setClauses = append(setClauses, "updated_at = now()")
	query := fmt.Sprintf(`
		update "user"
		set %s
		where user_id = $1
		returning user_id, username, coalesce(display_name, ''), email, password_hash, role, created_at, updated_at
	`, strings.Join(setClauses, ", "))

	row := r.pool.QueryRow(ctx, query, args...)
	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("%w: user not found", domain.ErrNotFound)
		}
		return domain.User{}, mapError(err, "user")
	}

	return user, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, userID int64) error {
	tag, err := r.pool.Exec(ctx, `delete from "user" where user_id = $1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: user not found", domain.ErrNotFound)
	}
	return nil
}

func scanUser(row interface{ Scan(dest ...any) error }) (domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}
