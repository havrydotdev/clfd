package file

import (
	"compress/gzip"
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	w := gzip.NewWriter(dst)
	defer w.Close()

	if _, err := io.Copy(w, src); err != nil {
		return err
	}

	return nil
}
