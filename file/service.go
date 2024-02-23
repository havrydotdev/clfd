package file

import (
	"context"

	"github.com/clfdrive/server/domain"
)

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
}

type Service struct {
	fileRepo FileRepository
}

func NewService(fileRepo FileRepository) *Service {
	return &Service{
		fileRepo,
	}
}

func (s *Service) Create(ctx context.Context, file *domain.File) error {
	return s.fileRepo.Create(ctx, file)
}
