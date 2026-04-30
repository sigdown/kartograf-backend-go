package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

type PostgresMapRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresMapRepository(pool *pgxpool.Pool) *PostgresMapRepository {
	return &PostgresMapRepository{pool: pool}
}

var _ usecase.MapRepository = (*PostgresMapRepository)(nil)

func (r *PostgresMapRepository) List(ctx context.Context) ([]domain.Map, error) {
	rows, err := r.pool.Query(ctx, `
		select uuid::text, created_by, coalesce(current_archive_id::text, ''), slug, title, coalesce(description, ''), coalesce(year, 0), created_at, updated_at
		from map
		order by created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps := make([]domain.Map, 0)
	for rows.Next() {
		m, err := scanMap(rows)
		if err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}

	return maps, rows.Err()
}

func (r *PostgresMapRepository) GetBySlug(ctx context.Context, slug string) (domain.Map, error) {
	return r.getOne(ctx, `
		select m.uuid::text, m.created_by, coalesce(m.current_archive_id::text, ''), m.slug, m.title, coalesce(m.description, ''), coalesce(m.year, 0), m.created_at, m.updated_at,
		       coalesce(a.archive_id::text, ''), coalesce(a.bucket, ''), coalesce(a.storage_key, ''), coalesce(a.uploaded_by, 0), coalesce(a.size_bytes, 0), coalesce(a.checksum, ''), coalesce(a.status, ''), coalesce(a.created_at, m.created_at), coalesce(a.updated_at, m.updated_at)
		from map m
		left join map_archive a on a.archive_id = m.current_archive_id
		where lower(m.slug) = lower($1)
	`, slug)
}

func (r *PostgresMapRepository) GetByID(ctx context.Context, mapID string) (domain.Map, error) {
	return r.getOne(ctx, `
		select m.uuid::text, m.created_by, coalesce(m.current_archive_id::text, ''), m.slug, m.title, coalesce(m.description, ''), coalesce(m.year, 0), m.created_at, m.updated_at,
		       coalesce(a.archive_id::text, ''), coalesce(a.bucket, ''), coalesce(a.storage_key, ''), coalesce(a.uploaded_by, 0), coalesce(a.size_bytes, 0), coalesce(a.checksum, ''), coalesce(a.status, ''), coalesce(a.created_at, m.created_at), coalesce(a.updated_at, m.updated_at)
		from map m
		left join map_archive a on a.archive_id = m.current_archive_id
		where m.uuid = $1
	`, mapID)
}

func (r *PostgresMapRepository) CreateWithArchive(ctx context.Context, m domain.Map, archive domain.MapArchive) (domain.Map, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Map{}, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		insert into map (uuid, created_by, slug, title, description, year)
		values ($1, $2, $3, $4, nullif($5, ''), nullif($6, 0))
	`, m.ID, m.CreatedBy, m.Slug, m.Title, m.Description, m.Year); err != nil {
		return domain.Map{}, mapError(err, "map")
	}

	if _, err := tx.Exec(ctx, `
		insert into map_archive (archive_id, map_id, bucket, storage_key, uploaded_by, size_bytes, checksum, status)
		values ($1, $2, $3, $4, $5, $6, nullif($7, ''), $8)
	`, archive.ID, archive.MapID, archive.Bucket, archive.StorageKey, archive.UploadedBy, archive.SizeBytes, archive.Checksum, archive.Status); err != nil {
		return domain.Map{}, mapError(err, "map archive")
	}

	if _, err := tx.Exec(ctx, `
		update map
		set current_archive_id = $2, updated_at = now()
		where uuid = $1
	`, m.ID, archive.ID); err != nil {
		return domain.Map{}, mapError(err, "map")
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Map{}, err
	}

	return r.GetByID(ctx, m.ID)
}

func (r *PostgresMapRepository) UpdateMetadata(ctx context.Context, mapID string, input usecase.UpdateMapMetadataInput) (domain.Map, error) {
	row := r.pool.QueryRow(ctx, `
		update map
		set title = $2,
		    description = nullif($3, ''),
		    year = nullif($4, 0),
		    updated_at = now()
		where uuid = $1
		returning uuid::text, created_by, coalesce(current_archive_id::text, ''), slug, title, coalesce(description, ''), coalesce(year, 0), created_at, updated_at
	`, mapID, input.Title, input.Description, input.Year)

	m, err := scanMap(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Map{}, fmt.Errorf("%w: map not found", domain.ErrNotFound)
		}
		return domain.Map{}, mapError(err, "map")
	}
	return m, nil
}

func (r *PostgresMapRepository) ReplaceArchive(ctx context.Context, mapID string, archive domain.MapArchive) (domain.MapArchive, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.MapArchive{}, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		update map_archive
		set bucket = $3,
		    storage_key = $4,
		    uploaded_by = $5,
		    size_bytes = $6,
		    checksum = nullif($7, ''),
		    status = $8,
		    updated_at = now()
		where map_id = $1 and archive_id = $2
	`, mapID, archive.ID, archive.Bucket, archive.StorageKey, archive.UploadedBy, archive.SizeBytes, archive.Checksum, domain.ArchiveStatusActive)
	if err != nil {
		return domain.MapArchive{}, mapError(err, "map archive")
	}
	if tag.RowsAffected() == 0 {
		return domain.MapArchive{}, fmt.Errorf("%w: archive not found", domain.ErrNotFound)
	}

	tag, err = tx.Exec(ctx, `
		update map
		set current_archive_id = $2, updated_at = now()
		where uuid = $1
	`, mapID, archive.ID)
	if err != nil {
		return domain.MapArchive{}, mapError(err, "map")
	}
	if tag.RowsAffected() == 0 {
		return domain.MapArchive{}, fmt.Errorf("%w: map not found", domain.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.MapArchive{}, err
	}

	return r.GetActiveArchive(ctx, mapID)
}

func (r *PostgresMapRepository) GetActiveArchive(ctx context.Context, mapID string) (domain.MapArchive, error) {
	row := r.pool.QueryRow(ctx, `
		select a.archive_id::text, a.map_id::text, a.bucket, a.storage_key, a.uploaded_by, coalesce(a.size_bytes, 0), coalesce(a.checksum, ''), a.status, a.created_at, a.updated_at
		from map m
		join map_archive a on a.archive_id = m.current_archive_id
		where m.uuid = $1
	`, mapID)

	archive, err := scanArchive(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.MapArchive{}, fmt.Errorf("%w: archive not found", domain.ErrNotFound)
		}
		return domain.MapArchive{}, err
	}
	return archive, nil
}

func (r *PostgresMapRepository) ListArchives(ctx context.Context, mapID string) ([]domain.MapArchive, error) {
	rows, err := r.pool.Query(ctx, `
		select archive_id::text, map_id::text, bucket, storage_key, uploaded_by, coalesce(size_bytes, 0), coalesce(checksum, ''), status, created_at, updated_at
		from map_archive
		where map_id = $1
		order by created_at desc
	`, mapID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	archives := make([]domain.MapArchive, 0)
	for rows.Next() {
		archive, err := scanArchive(rows)
		if err != nil {
			return nil, err
		}
		archives = append(archives, archive)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(archives) == 0 {
		row := r.pool.QueryRow(ctx, `select uuid::text, created_by, coalesce(current_archive_id::text, ''), slug, title, coalesce(description, ''), coalesce(year, 0), created_at, updated_at from map where uuid = $1`, mapID)
		if _, err := scanMap(row); err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: map not found", domain.ErrNotFound)
		}
	}

	return archives, nil
}

func (r *PostgresMapRepository) Delete(ctx context.Context, mapID string) error {
	tag, err := r.pool.Exec(ctx, `delete from map where uuid = $1`, mapID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: map not found", domain.ErrNotFound)
	}
	return nil
}

func (r *PostgresMapRepository) getOne(ctx context.Context, query string, arg string) (domain.Map, error) {
	row := r.pool.QueryRow(ctx, query, arg)
	m, archive, err := scanMapWithArchive(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Map{}, fmt.Errorf("%w: map not found", domain.ErrNotFound)
		}
		return domain.Map{}, err
	}

	if archive.ID != "" {
		m.CurrentArchive = &archive
	}
	return m, nil
}

func scanMap(row interface{ Scan(dest ...any) error }) (domain.Map, error) {
	var m domain.Map
	err := row.Scan(
		&m.ID,
		&m.CreatedBy,
		&m.CurrentArchiveID,
		&m.Slug,
		&m.Title,
		&m.Description,
		&m.Year,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}

func scanMapWithArchive(row interface{ Scan(dest ...any) error }) (domain.Map, domain.MapArchive, error) {
	var m domain.Map
	var archive domain.MapArchive
	err := row.Scan(
		&m.ID,
		&m.CreatedBy,
		&m.CurrentArchiveID,
		&m.Slug,
		&m.Title,
		&m.Description,
		&m.Year,
		&m.CreatedAt,
		&m.UpdatedAt,
		&archive.ID,
		&archive.Bucket,
		&archive.StorageKey,
		&archive.UploadedBy,
		&archive.SizeBytes,
		&archive.Checksum,
		&archive.Status,
		&archive.CreatedAt,
		&archive.UpdatedAt,
	)
	archive.MapID = m.ID
	return m, archive, err
}
