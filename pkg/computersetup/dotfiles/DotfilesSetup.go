package dotfiles

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
)

type DotfilesSetup struct {
	RepoURL      string
	LocalPath    string
	InstallScript string
}

func NewDotfilesSetup(repoURL, localPath, installScript string) *DotfilesSetup {
	return &DotfilesSetup{
		RepoURL:      repoURL,
		LocalPath:    localPath,
		InstallScript: installScript,
	}
}

func (d DotfilesSetup) Name() string {
	return "Dotfiles Setup"
}

func (d DotfilesSetup) Validate() (out *bytes.Buffer, err error) {
	out = bytes.NewBufferString("")
	
	// Check if dotfiles directory exists
	if _, statErr := os.Stat(d.LocalPath); os.IsNotExist(statErr) {
		_, _ = color.New(color.FgRed).Fprint(out, "NOT CLONED")
		out.WriteString(fmt.Sprintf(" %s", d.LocalPath))
		err = statErr
		return
	}
	
	// Check if it's a git repo
	gitDir := d.LocalPath + "/.git"
	if _, statErr := os.Stat(gitDir); os.IsNotExist(statErr) {
		_, _ = color.New(color.FgYellow).Fprint(out, "EXISTS BUT NOT A GIT REPO")
		out.WriteString(fmt.Sprintf(" %s", d.LocalPath))
		err = fmt.Errorf("directory exists but is not a git repository")
		return
	}
	
	_, _ = color.New(color.FgGreen).Fprint(out, "CLONED")
	out.WriteString(fmt.Sprintf(" %s", d.LocalPath))
	
	return
}

func (d DotfilesSetup) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	
	// Check if already exists
	validateOut, err := d.Validate()
	out.WriteString(validateOut.String() + "\n")
	if err == nil {
		out.WriteString("✅ Dotfiles already cloned, skipping clone\n")
	} else {
		// Clone the repository
		out.WriteString(fmt.Sprintf("Cloning %s to %s...\n", d.RepoURL, d.LocalPath))
		
		cmd := exec.Command("git", "clone", d.RepoURL, d.LocalPath)
		output, cloneErr := cmd.CombinedOutput()
		out.WriteString(string(output))
		
		if cloneErr != nil {
			_, _ = color.New(color.FgRed).Fprint(out, "❌ Failed to clone dotfiles\n")
			return out, cloneErr
		}
		
		_, _ = color.New(color.FgGreen).Fprint(out, "✅ Dotfiles cloned successfully\n")
	}
	
	// Always run install script if it exists (it's idempotent)
	if d.InstallScript != "" {
		scriptPath := d.LocalPath + "/" + d.InstallScript
		if _, statErr := os.Stat(scriptPath); os.IsNotExist(statErr) {
			out.WriteString(fmt.Sprintf("⚠️  Install script not found: %s\n", scriptPath))
		} else {
			out.WriteString(fmt.Sprintf("Running install script: %s\n", d.InstallScript))
			
			cmd := exec.Command("bash", scriptPath)
			cmd.Dir = d.LocalPath
			cmd.Env = os.Environ()
			
			output, scriptErr := cmd.CombinedOutput()
			out.WriteString(string(output))
			
			if scriptErr != nil {
				_, _ = color.New(color.FgRed).Fprint(out, "❌ Install script failed\n")
				return out, scriptErr
			}
			
			_, _ = color.New(color.FgGreen).Fprint(out, "✅ Install script completed\n")
		}
	}
	
	return out, nil
}

func (d DotfilesSetup) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	out.WriteString(fmt.Sprintf("❌ Cannot auto-uninstall dotfiles from %s - please remove manually\n", d.LocalPath))
	return out, nil
}
