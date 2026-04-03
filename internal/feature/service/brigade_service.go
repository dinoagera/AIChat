package service

import (
	"context"
	"log/slog"

	"github.com/dinoagera/AIChat/internal/domain"
)

type BrigadeService struct {
	log               *slog.Logger
	brigadeRepository BrigadeRepository
}

func NewBrigadeService(log *slog.Logger, brigadeRepository BrigadeRepository) *BrigadeService {
	return &BrigadeService{
		log:               log,
		brigadeRepository: brigadeRepository,
	}
}
func (b *BrigadeService) AddBrigade(ctx context.Context, name string, lat, lon float64, status string) error {
	if status == "" {
		status = "free"
	}
	exists, err := b.brigadeRepository.CheckName(ctx, name)
	if err != nil {
		b.log.Info("failed to check name", "err", err)
		return err
	}
	if exists {
		b.log.Info("name is exists")
		return domain.ErrBrigadeAlreadyExists
	}
	req := &domain.Brigade{
		Name:   name,
		Lat:    lat,
		Lon:    lon,
		Status: status,
	}
	if err := b.brigadeRepository.AddBrigade(ctx, req); err != nil {
		b.log.Info("failed to add brigade", "err", err)
		return err
	}
	return nil
}
func (b *BrigadeService) UpdateStatus(ctx context.Context, id int64, status string) error {
	exists, err := b.brigadeRepository.CheckBrigadeByID(ctx, id)
	if err != nil {
		b.log.Info("failed to check brigade by id", "err", err)
		return err
	}
	if !exists {
		return domain.ErrBridageNotExists
	}
	if err := b.brigadeRepository.UpdateStatus(ctx, id, status); err != nil {
		b.log.Info("failed to update status", "err", err)
		return err
	}
	return nil
}
