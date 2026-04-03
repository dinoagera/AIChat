package postgres

import (
	"context"

	domain "github.com/dinoagera/AIChat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrigadeRepository struct {
	pool *pgxpool.Pool
}

func NewBrigadesRepository(pool *pgxpool.Pool) *BrigadeRepository {
	return &BrigadeRepository{pool: pool}
}

func (br *BrigadeRepository) CheckName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := br.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM brigades WHERE name = $1)`, name).Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}
func (br *BrigadeRepository) AddBrigade(ctx context.Context, req *domain.Brigade) error {
	_, err := br.pool.Exec(ctx, `INSERT INTO brigades (name, lat, lon, status) VALUES($1,$2,$3,$4)`, req.Name, req.Lat, req.Lon, req.Status)
	if err != nil {
		return err
	}
	return nil
}
func (br *BrigadeRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := br.pool.Exec(ctx, `UPDATE brigades SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	return nil
}
func (br *BrigadeRepository) CheckBrigadeByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := br.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM brigades WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}
