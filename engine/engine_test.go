// engine/engine_test.go
package engine

// func TestStartEngine(t *testing.T) {
// 	ctx := context.Background()
// 	engineManager, err := NewEngineManager()
// 	if err != nil {
// 		t.Fatalf("Failed to create EngineManager: %v", err)
// 	}

// 	image := "hello-world:latest" // Use a simple image for testing
// 	containerName := "test_engine_container"
// 	cmd := []string{}
// 	ports := map[string]string{}
// 	resources := container.Resources{}

// 	// Clean up any existing container
// 	_ = engineManager.StopEngine(ctx, containerName)

// 	// Start the engine
// 	err = engineManager.StartEngine(ctx, image, containerName, cmd, ports, resources)
// 	if err != nil {
// 		t.Fatalf("Failed to start engine: %v", err)
// 	}

// 	// Verify the container is running
// 	containers, err := engineManager.ListEngines(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to list engines: %v", err)
// 	}

// 	found := false
// 	for _, c := range containers {
// 		for _, name := range c.Names {
// 			if name == "/"+containerName {
// 				found = true
// 				break
// 			}
// 		}
// 	}
// 	if !found {
// 		t.Errorf("Container %s not found after starting", containerName)
// 	}

// 	// Stop the engine
// 	err = engineManager.StopEngine(ctx, containerName)
// 	if err != nil {
// 		t.Errorf("Failed to stop engine: %v", err)
// 	}
// }
