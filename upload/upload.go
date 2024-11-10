package upload

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadZip saves the uploaded zip file to the specified upload directory after validating it contains index.html.
func UploadZip(file multipart.File, fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(data), fileHeader.Size)
	if err != nil {
		return "", err
	}
	hasIndex := false
	for _, f := range zipReader.File {
		if f.Name == "index.html" {
			hasIndex = true
			break
		}
	}
	if !hasIndex {
		return "", os.ErrInvalid
	}
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(uploadDir, fileHeader.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = dst.Write(data)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
