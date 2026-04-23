package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func mapError(err error, entity string) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s already exists", domain.ErrConflict, entity)
		case "23503":
			return fmt.Errorf("%w: invalid %s reference", domain.ErrInvalidInput, entity)
		case "23514":
			return fmt.Errorf("%w: invalid %s data", domain.ErrInvalidInput, entity)
		}
	}

	return err
}
