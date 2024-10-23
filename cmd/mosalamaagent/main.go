package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bodaay/mosalamaagent/engine"
	"github.com/bodaay/mosalamaagent/logging"
	"github.com/bodaay/mosalamaagent/model"
	"github.com/bodaay/mosalamaagent/monitoring"
)

func main() {
	// Initialize the logger
	logging.InitLogger()
	logging.Log.Info("Starting MosalamaAgent")

	// Initialize the engine manager
	engineManager, err := engine.NewEngineManager()
	if err != nil {
		logging.Log.Fatalf("Failed to initialize engine manager: %v", err)
	}

	// Initialize the model manager
	modelStorageDir := "/var/mosalamaagent/models"
	modelManager := model.NewModelManager(modelStorageDir)

	// Initialize the resource monitor
	resourceMonitor := monitoring.NewResourceMonitor()

	// Start resource monitoring
	stopChan := make(chan struct{})
	go resourceMonitor.StartMonitoring(30*time.Second, stopChan)

	// Simulate core functionality
	ctx := context.Background()

	// Example: Download a model
	modelURL := "https://models.example.com/your-model.bin"
	modelName := "your-model.bin"
	err = modelManager.DownloadModel(modelURL, modelName)
	if err != nil {
		logging.Log.Errorf("Failed to download model: %v", err)
	} else {
		logging.Log.Infof("Model downloaded successfully: %s", modelName)
	}

	// Example: Start an engine with the downloaded model
	image := "ghcr.io/engine/engine:latest" // Replace with the actual image name
	containerName := "mosalama_engine_container"
	cmd := []string{"--model", "/models/" + modelName}
	ports := map[string]string{"8000/tcp": "8000"} // Map container port to host port
	resources := engine.ContainerResources{
		CPUQuota: 200000,                 // Equivalent to 2 CPUs
		Memory:   4 * 1024 * 1024 * 1024, // 4GB RAM
		// Add GPU resources if necessary
	}

	err = engineManager.StartEngine(ctx, image, containerName, cmd, ports, resources)
	if err != nil {
		logging.Log.Errorf("Failed to start engine: %v", err)
	} else {
		logging.Log.Infof("Engine started successfully with container name: %s", containerName)
	}

	// Handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	logging.Log.Info("Shutdown signal received")

	// Stop the engine before exiting
	err = engineManager.StopEngine(ctx, containerName)
	if err != nil {
		logging.Log.Errorf("Failed to stop engine: %v", err)
	} else {
		logging.Log.Infof("Engine stopped successfully: %s", containerName)
	}

	close(stopChan)
	// Perform any cleanup if necessary
	logging.Log.Info("MosalamaAgent stopped")
}
