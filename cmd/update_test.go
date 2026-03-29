package cmd

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Command registration tests
// ---------------------------------------------------------------------------

func TestUpdateCmdRegistered(t *testing.T) {
	if UpdateCmd == nil {
		t.Fatal("UpdateCmd is nil")
	}
	if UpdateCmd.Use != "update" {
		t.Errorf("UpdateCmd.Use = %q, want %q", UpdateCmd.Use, "update")
	}
	if UpdateCmd.Short == "" {
		t.Error("UpdateCmd.Short is empty")
	}
}

func TestUpdateCmdAliases(t *testing.T) {
	aliases := UpdateCmd.Aliases
	want := map[string]bool{"up": false, "upgrade": false}

	for _, a := range aliases {
		if _, ok := want[a]; ok {
			want[a] = true
		}
	}

	for alias, found := range want {
		if !found {
			t.Errorf("UpdateCmd.Aliases = %v, want to contain %q", aliases, alias)
		}
	}
}

// ---------------------------------------------------------------------------
// Flag registration tests
// ---------------------------------------------------------------------------

func TestUpdateCmdFlags(t *testing.T) {
	flag := UpdateCmd.Flags().Lookup("dry-run")
	if flag == nil {
		t.Error("--dry-run flag not found on UpdateCmd")
	}

	flag = UpdateCmd.Flags().Lookup("managers")
	if flag == nil {
		t.Error("--managers flag not found on UpdateCmd")
	}
}

// ---------------------------------------------------------------------------
// Manager registry tests
// ---------------------------------------------------------------------------

func TestGetUpdatableManagers(t *testing.T) {
	managers := getUpdatableManagers()
	if len(managers) == 0 {
		t.Fatal("getUpdatableManagers() returned empty slice")
	}

	for _, mgr := range managers {
		if mgr.Name == "" {
			t.Error("manager has empty Name")
		}
		if mgr.Description == "" {
			t.Errorf("manager %q has empty Description", mgr.Name)
		}
		if len(mgr.Commands) == 0 {
			t.Errorf("manager %q has no Commands", mgr.Name)
		}
		for i, cmd := range mgr.Commands {
			if len(cmd) == 0 {
				t.Errorf("manager %q command[%d] is empty", mgr.Name, i)
			}
		}
	}
}

func TestGetUpdatableManagersExcludesDangerous(t *testing.T) {
	managers := getUpdatableManagers()
	excluded := []string{"pip", "cargo", "go"}

	nameSet := make(map[string]bool)
	for _, mgr := range managers {
		nameSet[mgr.Name] = true
	}

	for _, name := range excluded {
		if nameSet[name] {
			t.Errorf("getUpdatableManagers() should not include %q", name)
		}
	}
}

func TestSudoRequirements(t *testing.T) {
	managers := getUpdatableManagers()
	nameMap := make(map[string]packageManagerUpdate)
	for _, mgr := range managers {
		nameMap[mgr.Name] = mgr
	}

	wantSudo := []string{"apt", "snap", "dnf", "yum", "pacman"}
	for _, name := range wantSudo {
		mgr, ok := nameMap[name]
		if !ok {
			continue // manager not in registry (ok — it's there, but skip if not)
		}
		if !mgr.NeedsSudo {
			t.Errorf("manager %q should require sudo", name)
		}
	}

	noSudo := []string{"brew", "flatpak", "npm", "pipx", "gem", "choco", "winget"}
	for _, name := range noSudo {
		mgr, ok := nameMap[name]
		if !ok {
			continue
		}
		if mgr.NeedsSudo {
			t.Errorf("manager %q should NOT require sudo", name)
		}
	}
}

// ---------------------------------------------------------------------------
// Dry run smoke test
// ---------------------------------------------------------------------------

func TestRunUpdateDryRun(t *testing.T) {
	// Save and restore global flags
	oldDryRun := updateDryRun
	oldManagers := updateManagers
	defer func() {
		updateDryRun = oldDryRun
		updateManagers = oldManagers
	}()

	updateDryRun = true
	updateManagers = nil

	// Should not panic and should produce output
	output := captureOutput(func() {
		runUpdate()
	})

	if len(output) == 0 {
		t.Error("runUpdate() with --dry-run produced no output")
	}
	t.Logf("dry-run output length: %d", len(output))
}
