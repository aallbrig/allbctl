package osagnostic

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SSHKeyGitHubRegistration struct {
	KeyPath string
}

func NewSSHKeyGitHubRegistration() *SSHKeyGitHubRegistration {
	osSystem := NewOperatingSystem()
	defaultKeyPath := filepath.Join(osSystem.HomeDirectoryPath, ".ssh", "id_rsa.pub")
	return &SSHKeyGitHubRegistration{
		KeyPath: defaultKeyPath,
	}
}

func (s SSHKeyGitHubRegistration) Name() string {
	return "SSH Key GitHub Registration"
}

func (s SSHKeyGitHubRegistration) Validate() (out *bytes.Buffer, err error) {
	out = bytes.NewBufferString("")

	// Check if SSH key exists
	if _, statErr := os.Stat(s.KeyPath); os.IsNotExist(statErr) {
		_, _ = color.New(color.FgRed).Fprint(out, "SSH KEY NOT FOUND")
		out.WriteString(fmt.Sprintf(" %s", s.KeyPath))
		err = statErr
		return
	}

	// Check if GitHub CLI is available
	if _, ghErr := exec.LookPath("gh"); ghErr != nil {
		_, _ = color.New(color.FgRed).Fprint(out, "GITHUB CLI NOT FOUND")
		out.WriteString(" (required for SSH key registration)")
		err = ghErr
		return
	}

	// Check if authenticated with GitHub CLI
	cmd := exec.Command("gh", "auth", "status")
	if authErr := cmd.Run(); authErr != nil {
		_, _ = color.New(color.FgYellow).Fprint(out, "NOT AUTHENTICATED WITH GITHUB CLI")
		err = authErr
		return
	}

	// Check if SSH key is already registered
	keyContent, readErr := os.ReadFile(s.KeyPath)
	if readErr != nil {
		_, _ = color.New(color.FgRed).Fprint(out, "CANNOT READ SSH KEY")
		err = readErr
		return
	}

	// Extract just the key part (without comment)
	keyParts := strings.Fields(string(keyContent))
	if len(keyParts) < 2 {
		_, _ = color.New(color.FgRed).Fprint(out, "INVALID SSH KEY FORMAT")
		err = fmt.Errorf("invalid SSH key format")
		return
	}
	publicKey := keyParts[1]

	// List registered SSH keys
	cmd = exec.Command("gh", "ssh-key", "list")
	output, listErr := cmd.Output()
	if listErr != nil {
		_, _ = color.New(color.FgRed).Fprint(out, "CANNOT LIST GITHUB SSH KEYS")
		err = listErr
		return
	}

	if strings.Contains(string(output), publicKey) {
		_, _ = color.New(color.FgGreen).Fprint(out, "SSH KEY REGISTERED WITH GITHUB")
	} else {
		_, _ = color.New(color.FgYellow).Fprint(out, "SSH KEY NOT REGISTERED WITH GITHUB")
		err = fmt.Errorf("SSH key not registered")
	}

	return
}

func (s SSHKeyGitHubRegistration) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")

	// Generate SSH key if it doesn't exist
	if _, statErr := os.Stat(s.KeyPath); os.IsNotExist(statErr) {
		out.WriteString("Generating SSH key...\n")

		keyDir := filepath.Dir(s.KeyPath)
		if err := os.MkdirAll(keyDir, 0700); err != nil {
			_, _ = color.New(color.FgRed).Fprintf(out, "❌ Failed to create .ssh directory: %v\n", err)
			return out, err
		}

		privateKeyPath := strings.TrimSuffix(s.KeyPath, ".pub")
		cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", privateKeyPath, "-N", "")
		if keyGenErr := cmd.Run(); keyGenErr != nil {
			_, _ = color.New(color.FgRed).Fprintf(out, "❌ Failed to generate SSH key: %v\n", keyGenErr)
			return out, keyGenErr
		}

		_, _ = color.New(color.FgGreen).Fprint(out, "✅ SSH key generated\n")
	}

	// Check validation again
	validateOut, validationErr := s.Validate()
	if validateOut != nil {
		out.WriteString(validateOut.String() + "\n")
	}

	if validationErr != nil && strings.Contains(validationErr.Error(), "not registered") {
		// Register the SSH key with GitHub
		out.WriteString("Registering SSH key with GitHub...\n")

		// Get hostname for key title
		hostname, hostnameErr := os.Hostname()
		if hostnameErr != nil || hostname == "" {
			hostname = "allbctl-generated"
		}

		cmd := exec.Command("gh", "ssh-key", "add", s.KeyPath, "--title", fmt.Sprintf("allbctl-%s", hostname))
		output, addErr := cmd.CombinedOutput()
		out.WriteString(string(output))

		if addErr != nil {
			_, _ = color.New(color.FgRed).Fprintf(out, "❌ Failed to register SSH key: %v\n", addErr)
			return out, addErr
		}

		_, _ = color.New(color.FgGreen).Fprint(out, "✅ SSH key registered with GitHub\n")
	} else if validationErr != nil && strings.Contains(validationErr.Error(), "NOT AUTHENTICATED") {
		out.WriteString("⚠️  Please authenticate with GitHub CLI first:\n")
		out.WriteString("   gh auth login\n")
		return out, validationErr
	} else if validationErr != nil && strings.Contains(validationErr.Error(), "GITHUB CLI NOT FOUND") {
		out.WriteString("❌ GitHub CLI is required for SSH key registration\n")
		out.WriteString("Please install GitHub CLI first\n")
		return out, validationErr
	}

	return out, nil
}

func (s SSHKeyGitHubRegistration) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	out.WriteString("❌ Cannot auto-uninstall SSH key registration - please remove manually from GitHub\n")
	return out, nil
}
