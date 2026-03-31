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
func (b *BrigadeService) AddBrigade(ctx context.Context, req *domain.Brigade) error {
	if req.Status == "" {
		req.Status = "free"
	}
	if ok := b.brigadeRepository.CheckName(ctx, req.Name); ok {
		b.log.Info("name is exists")
		return domain.ErrBrigadeAlreadyExists
	}
	if err := b.brigadeRepository.AddBrigade(ctx, req); err != nil {
		b.log.Info("failed to add brigade", "err", err)
		return err
	}
	return nil
}
