package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SaveUploadedFile(r *http.Request, formFileName string) (string, error) {
	file, header, err := r.FormFile(formFileName)
	if err != nil {
		return "", nil // No file uploaded
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}
	return filePath, nil
}

func SaveUploadedFiles(r *http.Request, formFileName string) ([]string, error) {
	var filePaths []string
	files := r.MultipartForm.File[formFileName]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		os.MkdirAll("uploads", 0755)
		filePath := filepath.Join("uploads", fileHeader.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %v", err)
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			return nil, fmt.Errorf("error saving file: %v", err)
		}
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}
