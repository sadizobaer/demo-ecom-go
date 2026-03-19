package utilities

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func IsMethodValid(method string, validMethods []string) bool {
	for _, m := range validMethods {
		if method == m {
			return true
		}
	}
	return false
}

func IsFileSizeLimitExceeded(w http.ResponseWriter, r *http.Request, maxSize int64) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	err := r.ParseMultipartForm(maxSize)
	return err != nil
}

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, subPath string) (string, error) {
	filePath := "./uploads/" + subPath + "/" + header.Filename

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "Error creating directory", err
	}

	savedFile, err := os.Create(filePath)
	if err != nil {
		return "Error saving file", err
	}

	_, err = io.Copy(savedFile, file)
	if err != nil {
		return "Error copying file", err
	}

	defer file.Close()
	defer savedFile.Close()

	return filePath, nil
}
