package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintSystemInfo_Output(t *testing.T) {
	// Redirect stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	printSystemInfo()

	w.Close()
	os.Stdout = oldStdout

	var sb strings.Builder
	_, err = io.Copy(&sb, r)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := sb.String()

	// Check for expected sections (neofetch-style output)
	if !strings.Contains(output, "@") {
		t.Error("Output missing user@hostname header")
	}
	if !strings.Contains(output, "OS:") {
		t.Error("Output missing OS field")
	}
	if !strings.Contains(output, "Network:") {
		t.Error("Output missing Network section")
	}
	if !strings.Contains(output, "AI Agents:") {
		t.Error("Output missing AI Agents section")
	}
	if !strings.Contains(output, "Package Managers:") {
		t.Error("Output missing Package Managers section")
	}
	if !strings.Contains(output, "Packages:") {
		t.Error("Output missing Packages section")
	}
}

func Test_GetPackageManagerVersion(t *testing.T) {
	tests := []struct {
		name    string
		manager string
		wantErr bool
	}{
		{"npm", "npm", false},
		{"pip", "pip", false},
		{"unknown", "unknown-manager-xyz", false}, // Should not error, just return empty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := getPackageManagerVersion(tt.manager)
			// Version might be empty if manager not installed, that's ok
			t.Logf("Manager %s version: %s", tt.manager, version)
		})
	}
}

func Test_DetectAIAgents(t *testing.T) {
	agents := detectAIAgents()
	// May be empty if no AI agents installed, that's ok
	t.Logf("Detected AI agents: %v", agents)

	// Should not panic
	if len(agents) > 0 {
		for _, agent := range agents {
			if agent.Name == "" {
				t.Error("AI agent should have a name")
			}
		}
	}
}

func Test_GetAIAgentVersion(t *testing.T) {
	tests := []struct {
		name  string
		agent string
	}{
		{"copilot", "copilot"},
		{"claude", "claude"},
		{"ollama", "ollama"},
		{"unknown", "unknown-agent-xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := getAIAgentVersion(tt.agent)
			// Version might be empty if agent not installed
			t.Logf("Agent %s version: %s", tt.agent, version)
		})
	}
}

func Test_GetVersionManagerVersion(t *testing.T) {
	tests := []struct {
		name    string
		manager string
	}{
		{"nvm", "nvm"},
		{"pyenv", "pyenv"},
		{"rbenv", "rbenv"},
		{"rustup", "rustup"},
		{"unknown", "unknown-vm-xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := getVersionManagerVersion(tt.manager)
			// Version might be empty if manager not installed
			t.Logf("Version manager %s version: %s", tt.manager, version)
		})
	}
}

func Test_GetDetailedCPUInfo(t *testing.T) {
	cpuDetails := getDetailedCPUInfo()
	
	// Should always have some basic info
	if cpuDetails.ModelName == "" || cpuDetails.ModelName == "Unknown" {
		t.Error("CPU model name should not be empty or Unknown")
	}
	if cpuDetails.Architecture == "" {
		t.Error("CPU architecture should not be empty")
	}
	if cpuDetails.LogicalCores == 0 {
		t.Error("Logical cores should be greater than 0")
	}
	
	t.Logf("CPU Details: Model=%s, Arch=%s, Logical=%d, Physical=%d", 
		cpuDetails.ModelName, cpuDetails.Architecture, cpuDetails.LogicalCores, cpuDetails.PhysicalCores)
}

func Test_GetDetailedGPUInfo(t *testing.T) {
	gpus := getDetailedGPUInfo()
	
	// May be empty on systems without GPU, that's ok
	t.Logf("Detected %d GPU(s)", len(gpus))
	
	for i, gpu := range gpus {
		t.Logf("GPU %d: Name=%s, Vendor=%s", i, gpu.Name, gpu.Vendor)
		if gpu.Name == "" {
			t.Errorf("GPU %d should have a name", i)
		}
	}
}

func Test_DetectVendor(t *testing.T) {
	tests := []struct {
		name       string
		gpuName    string
		wantVendor string
	}{
		{"NVIDIA GeForce", "NVIDIA GeForce RTX 3080", "NVIDIA"},
		{"AMD Radeon", "AMD Radeon RX 6800", "AMD"},
		{"Intel UHD", "Intel UHD Graphics 630", "Intel"},
		{"Apple M1", "Apple M1", "Apple"},
		{"Microsoft Hyper-V", "Microsoft Corporation Hyper-V virtual VGA", "Microsoft"},
		{"ATI in word", "Corporation ATI Device", "AMD"}, // ATI as standalone word or part of compound
		{"Unknown", "Some Generic GPU", "Unknown"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vendor := detectVendor(tt.gpuName)
			if vendor != tt.wantVendor {
				t.Errorf("detectVendor(%s) = %s, want %s", tt.gpuName, vendor, tt.wantVendor)
			}
		})
	}
}

func Test_PrintCPUInfo(t *testing.T) {
	// Test that printCPUInfo doesn't panic
	details := CPUDetails{
		ModelName:     "Test CPU",
		Architecture:  "x86_64",
		LogicalCores:  8,
		PhysicalCores: 4,
		ThreadsPerCore: 2,
		BaseClock:     "3.5 GHz",
	}
	
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w
	
	printCPUInfo(details)
	
	w.Close()
	os.Stdout = oldStdout
	
	var sb strings.Builder
	_, err = io.Copy(&sb, r)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := sb.String()
	
	// Check for expected fields
	if !strings.Contains(output, "Test CPU") {
		t.Error("Output should contain CPU model name")
	}
	if !strings.Contains(output, "x86_64") {
		t.Error("Output should contain architecture")
	}
	if !strings.Contains(output, "3.5 GHz") {
		t.Error("Output should contain clock speed")
	}
}

func Test_PrintGPUInfo(t *testing.T) {
	// Test that printGPUInfo doesn't panic
	gpus := []GPUInfo{
		{
			Name:          "Test GPU",
			Vendor:        "Test Vendor",
			Memory:        "8 GB",
			Driver:        "1.0.0",
			ComputeCap:    "8.6",
			ClockGraphics: "1500 MHz",
			ClockMemory:   "7000 MHz",
		},
	}
	
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w
	
	printGPUInfo(gpus)
	
	w.Close()
	os.Stdout = oldStdout
	
	var sb strings.Builder
	_, err = io.Copy(&sb, r)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := sb.String()
	
	// Check for expected fields
	if !strings.Contains(output, "Test GPU") {
		t.Error("Output should contain GPU name")
	}
	if !strings.Contains(output, "Test Vendor") {
		t.Error("Output should contain vendor")
	}
	if !strings.Contains(output, "8 GB") {
		t.Error("Output should contain memory")
	}
}

