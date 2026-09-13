package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"kontursvet-api/internal/config"
	"kontursvet-api/internal/model"
	"kontursvet-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadService struct {
	repo       *repository.LeadRepository
	cfg        *config.Config
	uploadPath string
}

func NewLeadService(repo *repository.LeadRepository, cfg *config.Config) *LeadService {
	return &LeadService{
		repo:       repo,
		cfg:        cfg,
		uploadPath: cfg.Uploads.Path,
	}
}

func (s *LeadService) Create(ctx context.Context, lead *model.Lead) error {
	if err := s.repo.Create(ctx, lead); err != nil {
		return fmt.Errorf("lead service create: %w", err)
	}
	return nil
}

func (s *LeadService) SaveUploadFile(filename string, data []byte) (string, error) {
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return "", fmt.Errorf("lead service mkdir: %w", err)
	}

	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	path := filepath.Join(s.uploadPath, newFilename)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("lead service save file: %w", err)
	}

	return newFilename, nil
}

func (s *LeadService) List(ctx context.Context, limit, offset int) ([]model.Lead, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func GetDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
