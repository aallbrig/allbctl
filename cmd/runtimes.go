package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var RuntimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Detects programmer development environment runtimes and displays their versions.",
	Run: func(cmd *cobra.Command, args []string) {
		found := false
		for name, cmdArgs := range runtimeCommands {
			c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			output, err := c.CombinedOutput()
			if err == nil {
				fmt.Printf("%s: %s\n", name, parseVersion(string(output)))
				found = true
			}
		}
		if !found {
			fmt.Println("No known runtimes detected.")
		}
	},
}

var runtimeCommands = map[string][]string{
	"Node.js": {"node", "--version"},
	"Go":      {"go", "version"},
	"PHP":     {"php", "--version"},
	"Java":    {"java", "-version"},
	"Python":  {"python3", "--version"},
	// SQL Databases
	"MySQL":      {"mysql", "--version"},
	"PostgreSQL": {"psql", "--version"},
	"SQLite":     {"sqlite3", "--version"},
	"MariaDB":    {"mariadb", "--version"},
	"SQL Server": {"sqlcmd", "-?"},
	"Oracle":     {"sqlplus", "-version"},
	// NoSQL Databases
	"MongoDB":   {"mongod", "--version"},
	"Redis":     {"redis-server", "--version"},
	"Cassandra": {"cassandra", "-v"},
	// Kubernetes & Cloud
	"Kubernetes":       {"kubectl", "version", "--client", "--short"},
	"AWS CLI":          {"aws", "--version"},
	"Azure CLI":        {"az", "version"},
	"Google Cloud SDK": {"gcloud", "version"},
	// HashiCorp
	"Terraform": {"terraform", "version"},
	"Vault":     {"vault", "version"},
	"Consul":    {"consul", "version"},
	"Nomad":     {"nomad", "version"},
	// Package Managers
	"APT Packages": {"bash", "-c", "dpkg-query -f '${binary:Package}\n' -W | wc -l"},
}

func parseVersion(output string) string {
	// Simple version extraction: first line, trimmed
	return firstLine(output)
}

func firstLine(s string) string {
	for i, c := range s {
		if c == '\n' || c == '\r' {
			return s[:i]
		}
	}
	return s
}
