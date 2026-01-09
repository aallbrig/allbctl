package cmd

import (
	"testing"
)

func TestCheckPackageUpdates(t *testing.T) {
	tests := []struct {
		name    string
		manager string
		wantErr bool
	}{
		{"apt updates", "apt", false},
		{"flatpak updates", "flatpak", false},
		{"npm updates", "npm", false},
		{"pip updates", "pip", false},
		{"unknown manager", "unknown-manager", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := checkPackageUpdates(tt.manager)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPackageUpdates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Manager %s has %d updates available", tt.manager, count)
		})
	}
}

func TestFormatVersionWithUpdate(t *testing.T) {
	tests := []struct {
		name    string
		tool    string
		current string
		want    string
	}{
		{
			name:    "no update available",
			tool:    "npm",
			current: "11.6.2",
			want:    "11.6.2",
		},
		{
			name:    "unknown tool",
			tool:    "unknown-tool",
			current: "1.0.0",
			want:    "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatVersionWithUpdate(tt.tool, tt.current)
			// Since we don't have actual update data, it should return current version
			if got != tt.want {
				t.Errorf("formatVersionWithUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckOSUpdates(t *testing.T) {
	count, err := checkOSUpdates()
	if err != nil {
		t.Logf("checkOSUpdates() returned error (may be normal): %v", err)
	}
	t.Logf("OS has %d updates available", count)
}

func TestCheckAptUpdates(t *testing.T) {
	if !exists("apt") {
		t.Skip("apt not available on this system")
	}

	count, err := checkAptUpdates()
	if err != nil {
		t.Errorf("checkAptUpdates() error = %v", err)
		return
	}
	t.Logf("apt has %d updates available", count)
}

func TestCheckFlatpakUpdates(t *testing.T) {
	if !exists("flatpak") {
		t.Skip("flatpak not available on this system")
	}

	count, err := checkFlatpakUpdates()
	if err != nil {
		t.Errorf("checkFlatpakUpdates() error = %v", err)
		return
	}
	t.Logf("flatpak has %d updates available", count)
}

func TestCheckNpmUpdates(t *testing.T) {
	if !exists("npm") {
		t.Skip("npm not available on this system")
	}

	count, err := checkNpmUpdates()
	if err != nil {
		t.Errorf("checkNpmUpdates() error = %v", err)
		return
	}
	t.Logf("npm has %d global packages with updates available", count)
}

func TestCheckPipUpdates(t *testing.T) {
	if !exists("pip") && !exists("pip3") {
		t.Skip("pip not available on this system")
	}

	count, err := checkPipUpdates()
	if err != nil {
		t.Errorf("checkPipUpdates() error = %v", err)
		return
	}
	t.Logf("pip has %d packages with updates available", count)
}
