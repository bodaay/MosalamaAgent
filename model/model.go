package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ModelManager struct {
	StorageDir string
}

func NewModelManager(storageDir string) *ModelManager {
	return &ModelManager{
		StorageDir: storageDir,
	}
}
func (m *ModelManager) DownloadModel(url string, modelName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download model, status code: %d", resp.StatusCode)
	}

	modelPath := filepath.Join(m.StorageDir, modelName)
	out, err := os.Create(modelPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Optionally, add support for progress reporting and resuming
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Model downloaded: %s\n", modelPath)
	return nil
}
func (m *ModelManager) ListModels() ([]string, error) {
	files, err := os.ReadDir(m.StorageDir)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, file := range files {
		if !file.IsDir() {
			models = append(models, file.Name())
		}
	}
	return models, nil
}

func (m *ModelManager) DeleteModel(modelName string) error {
	modelPath := filepath.Join(m.StorageDir, modelName)
	if err := os.Remove(modelPath); err != nil {
		return err
	}
	fmt.Printf("Model deleted: %s\n", modelPath)
	return nil
}
