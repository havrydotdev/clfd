package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/internal/rest"
	"github.com/google/uuid"
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

func (s *Service) Create(ctx context.Context, upload *multipart.FileHeader, fileName, url string) (domain.File, error) {
	if fileName == "" {
		parts := strings.Split(upload.Filename, ".")

		fileName = fmt.Sprintf("%s.%s", uuid.NewString(), parts[len(parts) - 1])
	}
	
	_, err := SaveFile(upload, fileName)
	if err != nil {
		return domain.File{}, err
	}

	file := domain.File{
		Name: fileName,
		Location: fmt.Sprintf("%s/file/%s", url, fileName),
	}

	return file, s.fileRepo.Create(ctx, &file)
}

func (s *Service) GetFileName(ctx context.Context, fileName string) string {
	return path.Join(driveDir, fileName)
}
