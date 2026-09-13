package repository

import (
	"context"
	"fmt"
	"time"

	"kontursvet-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadRepository struct {
	pool *pgxpool.Pool	
}

func NewLeadRepository(pool *pgxpool.Pool) *LeadRepository {
	return &LeadRepository{pool: pool}
}

func (r *LeadRepository) Create(ctx context.Context, lead *model.Lead) error {
	query := `
		INSERT INTO leads (name, phone_digital, phone_format, home_type, location, message)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		lead.Name,
		lead.PhoneDigital,
		lead.PhoneFormat,
		lead.HomeType,
		lead.Location,
		lead.Message,
	).Scan(&lead.ID, &lead.CreatedAt, &lead.UpdatedAt)

	if err != nil {
		return fmt.Errorf("lead repository create: %w", err)
	}

	return nil
}

func (r *LeadRepository) List(ctx context.Context, limit, offset int) ([]model.Lead, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var leads []model.Lead

	countRows, err := r.pool.Query(ctx, "SELECT COUNT(*) FROM leads")
	if err != nil {
		return nil, 0, fmt.Errorf("lead repository count: %w", err)
	}
	defer countRows.Close()

	total := 0
	if countRows.Next() {
		countRows.Scan(&total)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, phone_digital, phone_format, home_type, location, message, created_at, updated_at
		FROM leads ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)

	if err != nil {
		return nil, 0, fmt.Errorf("lead repository list: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lead model.Lead
		if err := rows.Scan(
			&lead.ID, &lead.Name, &lead.PhoneDigital, &lead.PhoneFormat,
			&lead.HomeType, &lead.Location, &lead.Message,
			&lead.CreatedAt, &lead.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("lead repository scan: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, total, nil
}
