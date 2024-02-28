package pg

import (
	"context"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/file"
	"github.com/jackc/pgx/v5"
)

const (
	insertFileQuery = "INSERT INTO files (name, location, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
	deleteFileQuery = "DELETE FROM files WHERE name = $1 AND user_id = $2"
	findFilesByUser = "SELECT * FROM files WHERE user_id = $1"
)

type FileRepository struct {
	conn *pgx.Conn
}

func NewFileRepository(conn *pgx.Conn) file.FileRepository {
	return &FileRepository{
		conn,
	}
}

func (repo *FileRepository) Create(ctx context.Context, file *domain.File) (err error) {
	row := repo.conn.QueryRow(ctx, insertFileQuery, file.Name, file.Location, file.UserId)
	if err = row.Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt); err != nil {
		return
	}

	return nil
}

func (repo *FileRepository) FindByUser(ctx context.Context, userId int) ([]domain.File, error) {
	files := []domain.File{}
	rows, _ := repo.conn.Query(ctx, findFilesByUser, userId)
	for rows.Next() {
		var file domain.File
		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Location,
			&file.UpdatedAt,
			&file.CreatedAt,
			&file.UserId,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, rows.Err()
}

func (repo *FileRepository) Delete(ctx context.Context, fileName string, userId int) (err error) {
	row := repo.conn.QueryRow(ctx, deleteFileQuery, fileName, userId)
	if err = row.Scan(); err != nil && err != pgx.ErrNoRows {
		return
	}

	return nil
}
