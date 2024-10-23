Understood! Focusing on the core functionality of **MosalamaAgent** is a practical approach to get the essential features working before adding communication with the master node and other components. Let's proceed step by step to implement the core functionalities.

---

## **Implementation Plan for MosalamaAgent Core Functionality**

### **Overview**

We will focus on the following core components:

1. **Engine Management Module**: Managing language model serving engines using Docker.
2. **Model Management Module**: Handling model downloading and storage.
3. **Resource Monitoring Module**: Collecting system metrics (CPU, memory, GPU, etc.).
4. **Logging Module**: Implementing structured logging for debugging and monitoring.

---

### **Step 1: Set Up the Project Structure**

**a. Initialize a New Go Module**

Create a new directory for your project and initialize a Go module.

```bash
mkdir mosalamaagent
cd mosalamaagent
go mod init github.com/yourusername/mosalamaagent
```

**b. Create Project Directories**

Organize your code by creating packages for each module.

```
mosalamaagent/
├── cmd/
│   └── mosalamaagent/      // Main application entry point
├── engine/                 // Engine management
├── model/                  // Model management
├── monitoring/             // Resource monitoring
├── logging/                // Logging
├── utils/                  // Utility functions
└── go.mod
```

---

### **Step 2: Implement the Engine Management Module**

#### **a. Setup Docker Integration**

**Install Docker SDK for Go**

Add the Docker client library to your project.

```bash
go get github.com/docker/docker/client
```

**Create the Engine Package**

Create `engine/engine.go` for managing Docker containers.

```go
package engine

import (
    "context"
    "fmt"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
)

type EngineManager struct {
    dockerClient *client.Client
}

func NewEngineManager() (*EngineManager, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }
    return &EngineManager{
        dockerClient: cli,
    }, nil
}
```

#### **b. Implement Container Lifecycle Functions**

**Start an Engine**

```go
func (e *EngineManager) StartEngine(ctx context.Context, image string, containerName string, cmd []string, ports map[string]string, resources container.Resources) error {
    // Pull the image
    reader, err := e.dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
    if err != nil {
        return err
    }
    defer reader.Close()
    // Optionally, read the output from reader to monitor progress

    // Configure container
    containerConfig := &container.Config{
        Image: image,
        Cmd:   cmd,
    }

    hostConfig := &container.HostConfig{
        PortBindings:      natPortBindings(ports),
        Resources:         resources,
        RestartPolicy:     container.RestartPolicy{Name: "unless-stopped"},
    }

    resp, err := e.dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
    if err != nil {
        return err
    }

    // Start the container
    if err := e.dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        return err
    }

    fmt.Printf("Engine started with container ID: %s\n", resp.ID)
    return nil
}

// Helper function to convert port mappings
func natPortBindings(ports map[string]string) nat.PortMap {
    portMap := nat.PortMap{}
    for containerPort, hostPort := range ports {
        portMap[nat.Port(containerPort)] = []nat.PortBinding{
            {
                HostPort: hostPort,
            },
        }
    }
    return portMap
}
```

**Stop an Engine**

```go
func (e *EngineManager) StopEngine(ctx context.Context, containerName string) error {
    timeout := time.Second * 10
    if err := e.dockerClient.ContainerStop(ctx, containerName, &timeout); err != nil {
        return err
    }
    fmt.Printf("Engine stopped: %s\n", containerName)
    return nil
}
```

**List Running Engines**

```go
func (e *EngineManager) ListEngines(ctx context.Context) ([]types.Container, error) {
    containers, err := e.dockerClient.ContainerList(ctx, types.ContainerListOptions{})
    if err != nil {
        return nil, err
    }
    return containers, nil
}
```

#### **c. Handle Resource Allocation**

When starting containers, use the `Resources` field in `HostConfig` to set resource limits.

```go
resources := container.Resources{
    CPUQuota:  200000, // equivalent to 2 CPUs
    Memory:    4 * 1024 * 1024 * 1024, // 4GB RAM
    // For GPU resources, additional configurations are needed
}
```

---

### **Step 3: Implement the Model Management Module**

#### **a. Create the Model Package**

Create `model/model.go` to handle model downloading and storage.

```go
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
```

#### **b. Implement Model Downloading**

```go
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
```

#### **c. Manage Stored Models**

```go
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
```

---

### **Step 4: Implement the Resource Monitoring Module**

#### **a. Install Metrics Collection Libraries**

Add `gopsutil` for system metrics and `gonvml` for GPU metrics.

```bash
go get github.com/shirou/gopsutil/v3
go get github.com/mindprince/gonvml
```

#### **b. Create the Monitoring Package**

Create `monitoring/monitoring.go`.

```go
package monitoring

import (
    "fmt"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/mem"
    "time"
)

type ResourceMonitor struct {
}

func NewResourceMonitor() *ResourceMonitor {
    return &ResourceMonitor{}
}
```

#### **c. Implement System Metrics Collection**

**Collect CPU Usage**

```go
func (rm *ResourceMonitor) GetCPUUsage() (float64, error) {
    percentages, err := cpu.Percent(time.Second, false)
    if err != nil {
        return 0, err
    }
    if len(percentages) > 0 {
        return percentages[0], nil
    }
    return 0, fmt.Errorf("unable to get CPU usage")
}
```

**Collect Memory Usage**

```go
func (rm *ResourceMonitor) GetMemoryUsage() (uint64, uint64, float64, error) {
    vmStat, err := mem.VirtualMemory()
    if err != nil {
        return 0, 0, 0, err
    }
    return vmStat.Total, vmStat.Used, vmStat.UsedPercent, nil
}
```

**Collect Disk Usage**

```go
func (rm *ResourceMonitor) GetDiskUsage(path string) (uint64, uint64, float64, error) {
    usageStat, err := disk.Usage(path)
    if err != nil {
        return 0, 0, 0, err
    }
    return usageStat.Total, usageStat.Used, usageStat.UsedPercent, nil
}
```

**Collect GPU Usage (NVIDIA GPUs)**

```go
import (
    // ...
    nvml "github.com/mindprince/gonvml"
)

func (rm *ResourceMonitor) GetGPUUsage() ([]float64, error) {
    err := nvml.Initialize()
    if err != nil {
        return nil, err
    }
    defer nvml.Shutdown()

    deviceCount, err := nvml.DeviceCount()
    if err != nil {
        return nil, err
    }

    var gpuUsages []float64

    for i := uint(0); i < deviceCount; i++ {
        device, err := nvml.DeviceHandleByIndex(i)
        if err != nil {
            return nil, err
        }

        utilization, err := device.UtilizationRates()
        if err != nil {
            return nil, err
        }

        gpuUsages = append(gpuUsages, float64(utilization.GPU))
    }

    return gpuUsages, nil
}
```

#### **d. Periodic Monitoring**

Create a function to periodically collect and log metrics.

```go
func (rm *ResourceMonitor) StartMonitoring(interval time.Duration, stopChan <-chan struct{}) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            rm.collectAndLogMetrics()
        case <-stopChan:
            fmt.Println("Stopping monitoring")
            return
        }
    }
}

func (rm *ResourceMonitor) collectAndLogMetrics() {
    cpuUsage, err := rm.GetCPUUsage()
    if err != nil {
        fmt.Printf("Error collecting CPU usage: %v\n", err)
    } else {
        fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
    }
    // Similarly, collect memory, disk, and GPU usage
}
```

---

### **Step 5: Implement the Logging Module**

#### **a. Install a Logging Library**

Choose a structured logging library like `logrus` or `zap`.

```bash
go get github.com/sirupsen/logrus
```

#### **b. Create the Logging Package**

Create `logging/logger.go`.

```go
package logging

import (
    "os"

    "github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
    Log = logrus.New()
    Log.Out = os.Stdout
    Log.SetFormatter(&logrus.JSONFormatter{})
    // Set level based on configuration
    Log.SetLevel(logrus.InfoLevel)
}
```

#### **c. Use the Logger in Other Modules**

In other packages, import the logging package and use `logging.Log`.

```go
import (
    "github.com/yourusername/mosalamaagent/logging"
)

// Example usage
logging.Log.Info("Starting engine")
```

---

### **Step 6: Tie Everything Together in the Main Application**

Create `cmd/mosalamaagent/main.go`.

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourusername/mosalamaagent/engine"
    "github.com/yourusername/mosalamaagent/logging"
    "github.com/yourusername/mosalamaagent/model"
    "github.com/yourusername/mosalamaagent/monitoring"
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
    modelManager := model.NewModelManager("/var/mosalamaagent/models")

    // Initialize the resource monitor
    resourceMonitor := monitoring.NewResourceMonitor()

    // Start resource monitoring
    stopChan := make(chan struct{})
    go resourceMonitor.StartMonitoring(30*time.Second, stopChan)

    // Handle graceful shutdown
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
    <-signals
    logging.Log.Info("Shutdown signal received")

    close(stopChan)
    // Perform any cleanup if necessary
    logging.Log.Info("MosalamaAgent stopped")
}
```

---

### **Step 7: Testing Core Functionality**

#### **a. Testing Engine Management**

- Write a simple test to start and stop an engine.

```go
func testEngineManagement() {
    ctx := context.Background()
    engineManager, _ := engine.NewEngineManager()

    image := "my-engine-image:latest"
    containerName := "my_engine_container"
    cmd := []string{"--model", "/models/your-model"}
    ports := map[string]string{"8000/tcp": "8000"}
    resources := container.Resources{
        // Set resource limits
    }

    // Start Engine
    err := engineManager.StartEngine(ctx, image, containerName, cmd, ports, resources)
    if err != nil {
        logging.Log.Errorf("Failed to start engine: %v", err)
        return
    }

    // Perform operations...

    // Stop Engine
    err = engineManager.StopEngine(ctx, containerName)
    if err != nil {
        logging.Log.Errorf("Failed to stop engine: %v", err)
        return
    }
}
```

#### **b. Testing Model Management**

- Test downloading a model and listing stored models.

```go
func testModelManagement() {
    modelManager := model.NewModelManager("/var/mosalamaagent/models")
    url := "https://models.example.com/your-model.bin"
    modelName := "your-model.bin"

    err := modelManager.DownloadModel(url, modelName)
    if err != nil {
        logging.Log.Errorf("Failed to download model: %v", err)
        return
    }

    models, err := modelManager.ListModels()
    if err != nil {
        logging.Log.Errorf("Failed to list models: %v", err)
        return
    }

    logging.Log.Infof("Available models: %v", models)
}
```

#### **c. Testing Resource Monitoring**

- Ensure that resource metrics are being collected and logged.

---

### **Step 8: Configure and Run the Agent**

#### **a. Configuration File**

Create a configuration file (e.g., `config.yaml`) to store settings.

```yaml
# config.yaml
agent:
  storageDir: "/var/mosalamaagent/models"
  monitoringInterval: 30s
engine:
  defaultImage: "ghcr.io/engine/engine:latest"
  resources:
    cpuQuota: 200000
    memory: 4294967296
```

#### **b. Modify the Agent to Use Configurations**

Use a package like `viper` to read configurations.

```bash
go get github.com/spf13/viper
```

In `main.go`, initialize the configuration.

```go
import (
    // ...
    "github.com/spf13/viper"
)

func initConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("/etc/mosalamaagent/")
    viper.AddConfigPath(".")

    err := viper.ReadInConfig()
    if err != nil {
        logging.Log.Fatalf("Error reading config file: %v", err)
    }
}
```

Call `initConfig()` at the start of `main()`.

#### **c. Run the Agent**

Compile and run the agent.

```bash
go build -o mosalamaagent ./cmd/mosalamaagent
sudo ./mosalamaagent
```

Ensure that the agent starts correctly, logs output, and performs the core functions.

---

## **Next Steps**

- **Add Configuration Options**: Expand the configuration to cover all adjustable parameters.
- **Enhance Error Handling**: Improve error checks and handling for robustness.
- **Implement Unit Tests**: Write unit tests for each module to ensure reliability.
- **Prepare for Master Node Integration**: Design interfaces and placeholders for future integration with the master node.
- **Security Enhancements**: Even without the master node, start considering security in operations (e.g., validate inputs).
- **Documentation**: Document the code, functions, and usage instructions.

---

## **Conclusion**

By focusing on these core functionalities, you'll establish a solid foundation for **MosalamaAgent**. Once these components are working as expected, integrating communication with the master node and adding more advanced features will be more straightforward.

---

**Do you need further assistance with any of these steps, or is there a specific area you'd like to explore in more detail? I'm here to help you progress with your implementation.**