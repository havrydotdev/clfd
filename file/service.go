package file

import (
	"context"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/internal/rest"
)

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
}

type Service struct {
	fileRepo FileRepository
}

func NewService(fileRepo FileRepository) rest.FileService {
	return &Service{
		fileRepo,
	}
}

func (s *Service) Create(ctx context.Context, file *domain.File) error {
	return s.fileRepo.Create(ctx, file)
}
