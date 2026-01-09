package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var SecurityCmd = &cobra.Command{
	Use:   "security",
	Short: "Display security and authentication status",
	Long:  `Display information about SSH keys, GPG keys, and kernel keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintSecurityInfo()
	},
}

type SecurityInfo struct {
	SSHKeys     []string
	GPGKeys     []string
	KeyringInfo string
}

func PrintSecurityInfo() {
	fmt.Println("Security/Authentication Status:")
	fmt.Println()

	info := gatherSecurityInfo()

	// SSH Keys
	fmt.Printf("  SSH Keys (loaded in agent):\n")
	if len(info.SSHKeys) > 0 {
		for _, key := range info.SSHKeys {
			fmt.Printf("    - %s\n", key)
		}
	} else {
		fmt.Printf("    No keys loaded\n")
	}
	fmt.Println()

	// GPG Keys
	fmt.Printf("  GPG Keys:\n")
	if len(info.GPGKeys) > 0 {
		for _, key := range info.GPGKeys {
			fmt.Printf("    - %s\n", key)
		}
	} else {
		fmt.Printf("    No GPG keys found\n")
	}
	fmt.Println()

	// Kernel Keyring (Linux only)
	if runtime.GOOS == "linux" {
		fmt.Printf("  Kernel Keyring:\n")
		if info.KeyringInfo != "" {
			lines := strings.Split(info.KeyringInfo, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					fmt.Printf("    %s\n", line)
				}
			}
		} else {
			fmt.Printf("    No keys in user keyring\n")
		}
	}
}

func gatherSecurityInfo() *SecurityInfo {
	info := &SecurityInfo{
		SSHKeys: []string{},
		GPGKeys: []string{},
	}

	// Get SSH keys
	if exists("ssh-add") {
		cmd := exec.Command("ssh-add", "-l")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			for _, line := range lines {
				if line != "" && !strings.Contains(line, "has no identities") {
					// Parse ssh-add output: "2048 SHA256:... comment (RSA)"
					parts := strings.Fields(line)
					if len(parts) >= 3 {
						// Show key size, fingerprint (shortened), and comment
						fingerprint := parts[1]
						if len(fingerprint) > 20 {
							fingerprint = fingerprint[:20] + "..."
						}
						comment := ""
						if len(parts) > 2 {
							comment = strings.Join(parts[2:], " ")
						}
						info.SSHKeys = append(info.SSHKeys, fmt.Sprintf("%s %s %s", parts[0], fingerprint, comment))
					}
				}
			}
		}
	}

	// Get GPG keys
	if exists("gpg") {
		cmd := exec.Command("gpg", "--list-keys", "--keyid-format", "SHORT")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			var currentKey string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				// Look for pub lines: "pub   rsa4096/KEYID 2024-01-01 [SC]"
				if strings.HasPrefix(line, "pub ") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						currentKey = parts[1] // e.g., "rsa4096/KEYID"
					}
				}
				// Look for uid lines: "uid           [ultimate] Name <email>"
				if strings.HasPrefix(line, "uid ") && currentKey != "" {
					// Extract name/email
					uidParts := strings.SplitN(line, "]", 2)
					if len(uidParts) >= 2 {
						uid := strings.TrimSpace(uidParts[1])
						info.GPGKeys = append(info.GPGKeys, fmt.Sprintf("%s - %s", currentKey, uid))
						currentKey = "" // Reset for next key
					}
				}
			}
		}
	}

	// Get kernel keyring (Linux only)
	if runtime.GOOS == "linux" && exists("keyctl") {
		cmd := exec.Command("keyctl", "show", "@u")
		out, err := cmd.Output()
		if err == nil {
			output := strings.TrimSpace(string(out))
			// Count keys (lines that contain "keyring" or have key IDs)
			lines := strings.Split(output, "\n")
			keyCount := 0
			for _, line := range lines {
				if strings.Contains(line, ":") && !strings.Contains(line, "Keyring") {
					keyCount++
				}
			}
			if keyCount > 0 {
				info.KeyringInfo = fmt.Sprintf("%d keys in user keyring", keyCount)
			}
		}
	}

	return info
}
