// model/model_test.go
package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadModel(t *testing.T) {
	storageDir := "./test_models"
	modelManager := NewModelManager(storageDir)
	defer os.RemoveAll(storageDir) // Clean up after test

	// Create a test HTTP server to mock model download
	// Alternatively, if you have your own downloader, you can mock it

	modelContent := []byte("dummy model content")
	modelName := "test_model.bin"
	modelPath := filepath.Join(storageDir, modelName)

	// Simulate model download by writing to the file directly
	os.MkdirAll(storageDir, os.ModePerm)
	err := os.WriteFile(modelPath, modelContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write model file: %v", err)
	}

	// Verify the model exists
	models, err := modelManager.ListModels()
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}

	if len(models) == 0 {
		t.Fatalf("No models found in storage directory")
	}

	found := false
	for _, m := range models {
		if m == modelName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Model %s not found in list", modelName)
	}

	// Delete the model
	err = modelManager.DeleteModel(modelName)
	if err != nil {
		t.Errorf("Failed to delete model: %v", err)
	}

	// Verify the model is deleted
	models, err = modelManager.ListModels()
	if err != nil {
		t.Fatalf("Failed to list models after deletion: %v", err)
	}

	if len(models) != 0 {
		t.Errorf("Expected no models after deletion, found %d", len(models))
	}
}
