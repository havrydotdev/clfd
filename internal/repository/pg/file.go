package pg

import (
	"context"

	"github.com/clfdrive/server/domain"
	"github.com/jackc/pgx/v5"
)

const (
	insertQuery = "INSERT INTO files (name, location) VALUES ($1, $2) RETURNING id, created_at, updated_at"
)

type FileRepository struct {
	conn *pgx.Conn
}

func NewFileRepository(conn *pgx.Conn) *FileRepository {
	return &FileRepository{
		conn,
	}
}

func (repo *FileRepository) Create(ctx context.Context, file *domain.File) (err error) {
	row := repo.conn.QueryRow(ctx, insertQuery, file.Name, file.Location)
	if err = row.Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt); err != nil {
		return
	}

	return nil
}
