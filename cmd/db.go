package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dbDetailFlag bool
)

// DbCmd represents the db command
var DbCmd = &cobra.Command{
	Use:   "db [database]",
	Short: "Display detected databases and their information",
	Long: `Display detected database management systems (DBMS) and their information.

Shows database clients, versions, running status, and database files.

Examples:
  allbctl status db                  # Show all detected databases
  allbctl status db sqlite3          # Show only SQLite3 info
  allbctl status db postgres         # Show only PostgreSQL info
  allbctl status db --detail         # Show detailed info for all databases
  allbctl status db sqlite3 --detail # Show detailed SQLite3 info with .db files`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// Show specific database
			dbName := args[0]
			showDatabaseInfo(dbName, dbDetailFlag)
		} else {
			// Show all databases
			showAllDatabases(dbDetailFlag)
		}
	},
}

func init() {
	DbCmd.Flags().BoolVarP(&dbDetailFlag, "detail", "d", false, "Show detailed information including database files and environment variables")
}

type DatabaseInfo struct {
	Name          string
	ClientBinary  string
	ServerBinary  string
	ClientVersion string
	ServerVersion string
	IsRunning     bool
	DatabaseFiles []string
	EnvVars       map[string]string
	OtherBinaries []string
}

// Database configurations
var databaseConfigs = map[string]struct {
	clientBinary   string
	serverBinary   string
	versionArgs    []string
	serverCheck    []string
	fileExtensions []string
	envVarPrefixes []string
}{
	"sqlite3": {
		clientBinary:   "sqlite3",
		versionArgs:    []string{"--version"},
		fileExtensions: []string{".db", ".sqlite", ".sqlite3"},
		envVarPrefixes: []string{"SQLITE_"},
	},
	"mysql": {
		clientBinary:   "mysql",
		serverBinary:   "mysqld",
		versionArgs:    []string{"--version"},
		serverCheck:    []string{"pgrep", "-x", "mysqld"},
		envVarPrefixes: []string{"MYSQL_"},
	},
	"mariadb": {
		clientBinary:   "mariadb",
		serverBinary:   "mariadbd",
		versionArgs:    []string{"--version"},
		serverCheck:    []string{"pgrep", "-x", "mariadbd"},
		envVarPrefixes: []string{"MYSQL_", "MARIADB_"},
	},
	"postgres": {
		clientBinary:   "psql",
		serverBinary:   "postgres",
		versionArgs:    []string{"--version"},
		serverCheck:    []string{"pgrep", "-x", "postgres"},
		envVarPrefixes: []string{"PG", "POSTGRES_"},
	},
	"mongodb": {
		clientBinary:   "mongosh",
		serverBinary:   "mongod",
		versionArgs:    []string{"--version"},
		serverCheck:    []string{"pgrep", "-x", "mongod"},
		envVarPrefixes: []string{"MONGO_"},
	},
	"redis": {
		clientBinary:   "redis-cli",
		serverBinary:   "redis-server",
		versionArgs:    []string{"--version"},
		serverCheck:    []string{"pgrep", "-x", "redis-server"},
		envVarPrefixes: []string{"REDIS_"},
	},
	"cassandra": {
		clientBinary:   "cqlsh",
		serverBinary:   "cassandra",
		versionArgs:    []string{"--version"},
		serverCheck:    []string{"pgrep", "-f", "cassandra"},
		envVarPrefixes: []string{"CASSANDRA_"},
	},
	"oracle": {
		clientBinary:   "sqlplus",
		versionArgs:    []string{"-version"},
		envVarPrefixes: []string{"ORACLE_", "TNS_"},
	},
	"sqlserver": {
		clientBinary:   "sqlcmd",
		versionArgs:    []string{"-?"},
		envVarPrefixes: []string{"MSSQL_"},
	},
}

func detectDatabase(dbName string) *DatabaseInfo {
	config, exists := databaseConfigs[dbName]
	if !exists {
		return nil
	}

	info := &DatabaseInfo{
		Name:          dbName,
		ClientBinary:  config.clientBinary,
		ServerBinary:  config.serverBinary,
		EnvVars:       make(map[string]string),
		OtherBinaries: []string{},
	}

	// Check client binary
	if clientPath, err := exec.LookPath(config.clientBinary); err == nil {
		info.ClientBinary = clientPath
		// Get version
		if len(config.versionArgs) > 0 {
			cmd := exec.Command(config.clientBinary, config.versionArgs...)
			if output, err := cmd.CombinedOutput(); err == nil {
				info.ClientVersion = strings.TrimSpace(string(output))
			}
		}
	} else {
		return nil // Client not installed
	}

	// Check server binary
	if config.serverBinary != "" {
		if serverPath, err := exec.LookPath(config.serverBinary); err == nil {
			info.ServerBinary = serverPath
			// Get server version (often same as client)
			if len(config.versionArgs) > 0 {
				cmd := exec.Command(config.serverBinary, config.versionArgs...)
				if output, err := cmd.CombinedOutput(); err == nil {
					info.ServerVersion = strings.TrimSpace(string(output))
				}
			}
		}

		// Check if server is running
		if len(config.serverCheck) > 0 {
			cmd := exec.Command(config.serverCheck[0], config.serverCheck[1:]...)
			if err := cmd.Run(); err == nil {
				info.IsRunning = true
			}
		}
	}

	// Find database files (only for file-based databases)
	if len(config.fileExtensions) > 0 {
		info.DatabaseFiles = findDatabaseFiles(config.fileExtensions)
	}

	// Collect environment variables
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			for _, prefix := range config.envVarPrefixes {
				if strings.HasPrefix(key, prefix) {
					info.EnvVars[key] = value
					break
				}
			}
		}
	}

	return info
}

func findDatabaseFiles(extensions []string) []string {
	var files []string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return files
	}

	srcDir := filepath.Join(homeDir, "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return files
	}

	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			for _, ext := range extensions {
				if strings.HasSuffix(strings.ToLower(info.Name()), ext) {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	return files
}

func showDatabaseInfo(dbName string, detailed bool) {
	info := detectDatabase(dbName)
	if info == nil {
		fmt.Printf("Database '%s' not detected on this system\n", dbName)
		os.Exit(1)
	}

	printDatabaseInfo(info, detailed)
}

func showAllDatabases(detailed bool) {
	var detectedDatabases []string
	var allInfo []*DatabaseInfo

	// Detect all databases
	for dbName := range databaseConfigs {
		if info := detectDatabase(dbName); info != nil {
			detectedDatabases = append(detectedDatabases, dbName)
			allInfo = append(allInfo, info)
		}
	}

	if len(detectedDatabases) == 0 {
		fmt.Println("No databases detected")
		return
	}

	// Print summary or detailed info
	if detailed {
		for i, info := range allInfo {
			if i > 0 {
				fmt.Println()
			}
			printDatabaseInfo(info, true)
		}
		// Print summary at the end
		fmt.Printf("\nDatabases detected: %s\n", strings.Join(detectedDatabases, ", "))
	} else {
		for _, info := range allInfo {
			printDatabaseSummary(info)
		}
	}
}

func printDatabaseSummary(info *DatabaseInfo) {
	status := "installed"
	if info.IsRunning {
		status = "running"
	}

	version := info.ClientVersion
	if version == "" {
		version = "unknown"
	} else {
		// Shorten version for summary
		lines := strings.Split(version, "\n")
		if len(lines) > 0 {
			version = strings.TrimSpace(lines[0])
			// Truncate if too long
			if len(version) > 60 {
				version = version[:57] + "..."
			}
		}
	}

	fileCount := ""
	if len(info.DatabaseFiles) > 0 {
		fileCount = fmt.Sprintf(", %d .db files", len(info.DatabaseFiles))
	}

	fmt.Printf("%-12s %-10s %s%s\n", info.Name+":", status, version, fileCount)
}

func printDatabaseInfo(info *DatabaseInfo, detailed bool) {
	fmt.Printf("%s:\n", strings.Title(info.Name))
	fmt.Println("----------------------------------------")

	// Client information
	fmt.Printf("  Client Binary:  %s\n", info.ClientBinary)
	if info.ClientVersion != "" {
		fmt.Printf("  Client Version: %s\n", info.ClientVersion)
	}

	// Server information
	if info.ServerBinary != "" {
		fmt.Printf("  Server Binary:  %s\n", info.ServerBinary)
		if info.ServerVersion != "" && info.ServerVersion != info.ClientVersion {
			fmt.Printf("  Server Version: %s\n", info.ServerVersion)
		}
		if info.IsRunning {
			fmt.Printf("  Status:         RUNNING\n")
		} else {
			fmt.Printf("  Status:         not running\n")
		}
	}

	if detailed {
		// Database files
		if len(info.DatabaseFiles) > 0 {
			fmt.Printf("  Database Files: (%d found in ~/src)\n", len(info.DatabaseFiles))
			for _, file := range info.DatabaseFiles {
				// Convert to relative path from home
				homeDir, _ := os.UserHomeDir()
				relPath := strings.Replace(file, homeDir, "~", 1)
				fmt.Printf("    - %s\n", relPath)
			}
		}

		// Environment variables
		if len(info.EnvVars) > 0 {
			fmt.Printf("  Environment Variables:\n")
			for key, value := range info.EnvVars {
				fmt.Printf("    %s=%s\n", key, value)
			}
		}
	} else if len(info.DatabaseFiles) > 0 {
		fmt.Printf("  Database Files: %d found in ~/src\n", len(info.DatabaseFiles))
	}
}

// PrintDatabaseSummaryForStatus prints a one-line summary for the main status command
func PrintDatabaseSummaryForStatus() {
	var detected []string

	for dbName := range databaseConfigs {
		if info := detectDatabase(dbName); info != nil {
			status := ""
			if info.IsRunning {
				status = " (running)"
			}
			detected = append(detected, dbName+status)
		}
	}

	if len(detected) > 0 {
		fmt.Printf("Databases: %s\n", strings.Join(detected, ", "))
	}
}
