package monitoring

import (
	"testing"
)

func TestGetCPUUsage(t *testing.T) {
	rm := NewResourceMonitor()
	usage, err := rm.GetCPUUsage()
	if err != nil {
		t.Fatalf("Failed to get CPU usage: %v", err)
	}
	t.Logf("CPU Usage: %.2f%%", usage)
}

func TestGetMemoryUsage(t *testing.T) {
	rm := NewResourceMonitor()
	total, used, percent, err := rm.GetMemoryUsage()
	if err != nil {
		t.Fatalf("Failed to get memory usage: %v", err)
	}
	t.Logf("Memory Usage: Total=%d, Used=%d, UsedPercent=%.2f%%", total, used, percent)
}
