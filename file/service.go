package file

import (
	"context"
	"fmt"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"

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
	parts := strings.Split(upload.Filename, ".")
	fileName := fmt.Sprintf("%s.%s", uuid.NewString(), parts[len(parts)-1])

	filePath := s.GetFileName(ctx, fileName, userId)
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
	}

	return file, s.fileRepo.Create(ctx, &file)
}

func (s *Service) GetFileName(ctx context.Context, fileName string, userId int) string {
	return path.Join(
		driveDir,
		strconv.Itoa(userId),
		fileName,
	)
}
