package file

import (
	"context"
	"fmt"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/internal/rest"
	"github.com/google/uuid"
)

var (
	currDir, _ = os.Getwd()
	driveDir   = path.Join(currDir, ".drive")
)

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	FindByUser(ctx context.Context, userId int) ([]domain.File, error)
	Delete(ctx context.Context, fileName string, userId int) error
}

type Service struct {
	fileRepo FileRepository
}

func NewService(fileRepo FileRepository) rest.FileService {
	return &Service{
		fileRepo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	upload *multipart.FileHeader,
	url string,
	userId int,
) (domain.File, error) {
	ext := filepath.Ext(upload.Filename)
	fileName := uuid.NewString() + ext
	filePath := s.GetFilePath(ctx, fileName, userId)

	dirPath := path.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, fs.ModePerm)
	}

	err := SaveFile(upload, filePath)
	if err != nil {
		return domain.File{}, err
	}

	file := domain.File{
		Name:     fileName,
		Location: fmt.Sprintf("%s/file/%s", url, fileName),
		UserId:   userId,
	}

	return file, s.fileRepo.Create(ctx, &file)
}

func (s *Service) Delete(ctx context.Context, fileName string, userId int) error {
	filePath := s.GetFilePath(ctx, fileName, userId)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return s.fileRepo.Delete(ctx, fileName, userId)
}

func (s *Service) FindByUser(ctx context.Context, userId int) ([]domain.File, error) {
	return s.fileRepo.FindByUser(ctx, userId)
}

func (s *Service) GetFilePath(ctx context.Context, fileName string, userId int) string {
	return path.Join(
		driveDir,
		strconv.Itoa(userId),
		fileName+".gz",
	)
}

func (s *Service) ReadGzip(ctx context.Context, filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
