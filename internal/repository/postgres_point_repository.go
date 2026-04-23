package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

type PostgresPointRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPointRepository(pool *pgxpool.Pool) *PostgresPointRepository {
	return &PostgresPointRepository{pool: pool}
}

var _ usecase.PointRepository = (*PostgresPointRepository)(nil)

func (r *PostgresPointRepository) ListByOwner(ctx context.Context, ownerID int64) ([]domain.Point, error) {
	rows, err := r.pool.Query(ctx, `
		select point_id, owner_id, name, coalesce(description, ''), lat, lon, created_at, updated_at
		from point
		where owner_id = $1
		order by created_at desc
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]domain.Point, 0)
	for rows.Next() {
		point, err := scanPoint(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, rows.Err()
}

func (r *PostgresPointRepository) Create(ctx context.Context, point domain.Point) (domain.Point, error) {
	row := r.pool.QueryRow(ctx, `
		insert into point (owner_id, name, description, lat, lon)
		values ($1, $2, nullif($3, ''), $4, $5)
		returning point_id, owner_id, name, coalesce(description, ''), lat, lon, created_at, updated_at
	`, point.OwnerID, point.Name, point.Description, point.Lat, point.Lon)

	created, err := scanPoint(row)
	if err != nil {
		return domain.Point{}, mapError(err, "point")
	}
	return created, nil
}

func (r *PostgresPointRepository) GetByID(ctx context.Context, pointID int64) (domain.Point, error) {
	row := r.pool.QueryRow(ctx, `
		select point_id, owner_id, name, coalesce(description, ''), lat, lon, created_at, updated_at
		from point
		where point_id = $1
	`, pointID)

	point, err := scanPoint(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Point{}, fmt.Errorf("%w: point not found", domain.ErrNotFound)
		}
		return domain.Point{}, err
	}
	return point, nil
}

func (r *PostgresPointRepository) Update(ctx context.Context, pointID int64, input usecase.UpdatePointInput) (domain.Point, error) {
	row := r.pool.QueryRow(ctx, `
		update point
		set name = $2,
		    description = nullif($3, ''),
		    lat = $4,
		    lon = $5,
		    updated_at = now()
		where point_id = $1
		returning point_id, owner_id, name, coalesce(description, ''), lat, lon, created_at, updated_at
	`, pointID, input.Name, input.Description, input.Lat, input.Lon)

	point, err := scanPoint(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Point{}, fmt.Errorf("%w: point not found", domain.ErrNotFound)
		}
		return domain.Point{}, mapError(err, "point")
	}
	return point, nil
}

func (r *PostgresPointRepository) Delete(ctx context.Context, pointID int64) error {
	tag, err := r.pool.Exec(ctx, `delete from point where point_id = $1`, pointID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: point not found", domain.ErrNotFound)
	}
	return nil
}

func scanPoint(row interface{ Scan(dest ...any) error }) (domain.Point, error) {
	var point domain.Point
	err := row.Scan(
		&point.ID,
		&point.OwnerID,
		&point.Name,
		&point.Description,
		&point.Lat,
		&point.Lon,
		&point.CreatedAt,
		&point.UpdatedAt,
	)
	return point, err
}
