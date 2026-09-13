package repository

import (
	"context"
	"fmt"
	"time"

	"kontursvet-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PortfolioRepository struct {
	pool *pgxpool.Pool
}

func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{pool: pool}
}

func (r *PortfolioRepository) ListCases(ctx context.Context) ([]model.PortfolioCase, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, case, title, description, task, location, term, team, period, features, created_at, updated_at
		FROM portfolio_cases ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("portfolio repository list cases: %w", err)
	}
	defer rows.Close()

	var cases []model.PortfolioCase

	for rows.Next() {
		var case_item model.PortfolioCase
		if err := rows.Scan(
			&case_item.ID, &case_item.Name, &case_item.Case, &case_item.Title, &case_item.Description,
			&case_item.Task, &case_item.Location, &case_item.Term, &case_item.Team, &case_item.Period,
			&case_item.Features, &case_item.CreatedAt, &case_item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("portfolio repository scan case: %w", err)
		}
		cases = append(cases, case_item)
	}

	return cases, nil
}

func (r *PortfolioRepository) GetCaseWithPhotos(ctx context.Context, id int64) (*model.PortfolioCase, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var case_item model.PortfolioCase
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, case, title, description, task, location, term, team, period, features, created_at, updated_at
		FROM portfolio_cases WHERE id = $1`, id).Scan(
		&case_item.ID, &case_item.Name, &case_item.Case, &case_item.Title, &case_item.Description,
		&case_item.Task, &case_item.Location, &case_item.Term, &case_item.Team, &case_item.Period,
		&case_item.Features, &case_item.CreatedAt, &case_item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("portfolio repository get case: %w", err)
	}

	// Get works
	worksRows, err := r.pool.Query(ctx,
		"SELECT work FROM portfolio_case_works WHERE case_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("portfolio repository get works: %w", err)
	}
	defer worksRows.Close()

	for worksRows.Next() {
		var work string
		worksRows.Scan(&work)
		case_item.Works = append(case_item.Works, work)
	}

	// Get photos
	photosRows, err := r.pool.Query(ctx,
		"SELECT gallery_key, photo_path, alt, caption FROM portfolio_photos WHERE case_id = $1 ORDER BY gallery_key", id)
	if err != nil {
		return nil, fmt.Errorf("portfolio repository get photos: %w", err)
	}
	defer photosRows.Close()

	for photosRows.Next() {
		var path, alt, caption string
		var key int
		photosRows.Scan(&key, &path, &alt, &caption)
		// Can't directly append to nested struct, skip for now
		// Will add photos field to model later
		_ = key
		_ = path
		_ = alt
		_ = caption
	}

	return &case_item, nil
}