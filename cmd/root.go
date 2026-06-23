package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/aallbrig/allbctl/pkg/telemetry"
)

var cfgFile string
var debugMode bool

// telemetryShutdown is set by initTelemetry and called in postRunTelemetry.
var telemetryShutdown func(context.Context) error

// commandStartTime records when the current command started for duration metrics.
var commandStartTime time.Time

var rootCmd = &cobra.Command{
	Use:   "allbctl",
	Short: "allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks",
	Long: `allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks.

Example commands for allbctl:

$ allbctl bootstrap status
$ allbctl bootstrap install
$ allbctl status
$ allbctl status runtimes              # Show detected programming runtimes
$ allbctl status projects              # Show git repositories in ~/src
$ allbctl status list-packages         # Show package counts from all package managers
$ allbctl status db                    # Show detected databases and their status
$ allbctl status network               # Show network interface information
$ allbctl status containers            # Show container/virtualization info
$ allbctl status security              # Show SSH keys, GPG keys, and keyring
$ allbctl status systemctl             # Show systemd service counts
$ allbctl status git                   # Show git global configuration
$ allbctl status ports                 # Show listening ports
$ allbctl status cloud-native          # Show cloud CLI versions and profiles
$ allbctl status cloud-native aws      # Show detailed AWS resources by region
$ allbctl update                       # Update all detected package managers
$ allbctl update --dry-run             # Preview updates without executing
$ allbctl update --managers apt,npm    # Only update apt and npm
`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initTelemetry(cmd, args)
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return postRunTelemetry(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() //nolint:errcheck // Help errors are not critical
	},
}

// Execute comment for execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(BootstrapCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(UpdateCmd)

	// Add subcommands to status
	StatusCmd.AddCommand(RuntimesCmd)
	StatusCmd.AddCommand(ListPackagesCmd)
	StatusCmd.AddCommand(ProjectsCmd)
	StatusCmd.AddCommand(DbCmd)
	StatusCmd.AddCommand(NetworkCmd)
	StatusCmd.AddCommand(ContainersCmd)
	StatusCmd.AddCommand(SecurityCmd)
	StatusCmd.AddCommand(SystemctlCmd)
	StatusCmd.AddCommand(GitConfigCmd)
	StatusCmd.AddCommand(PortsCmd)
	StatusCmd.AddCommand(CloudNativeCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.allbctl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug telemetry (structured logs, traces, and metrics to stderr)")

	rootCmd.SetVersionTemplate(fmt.Sprintf("allbctl %s (commit %s)\n", Version, Commit))
}

// initTelemetry is called by PersistentPreRunE on every command. It sets up
// the OTel providers and starts a root span for the command invocation.
func initTelemetry(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	shutdown, err := telemetry.Setup(ctx, debugMode)
	if err != nil {
		// Non-fatal: log to stderr and continue without telemetry
		fmt.Fprintf(os.Stderr, "telemetry setup warning: %v\n", err)
		telemetryShutdown = func(context.Context) error { return nil }
		return nil
	}
	telemetryShutdown = shutdown
	commandStartTime = time.Now()

	// Start a root span that covers the entire command execution.
	// The span is stored in ctx and ended in postRunTelemetry.
	ctx, _ = otel.Tracer("github.com/aallbrig/allbctl").Start(
		ctx, cmd.CommandPath(),
		trace.WithAttributes(
			attribute.String("command", cmd.CommandPath()),
			attribute.StringSlice("args", args),
			attribute.String("version", Version),
		),
	)
	cmd.SetContext(ctx)

	telemetry.Logger.InfoContext(ctx, "command.start",
		"command", cmd.CommandPath(),
		"args", args,
		"debug", debugMode,
		"version", Version,
	)

	return nil
}

// postRunTelemetry ends the command span, records metrics, and flushes providers.
func postRunTelemetry(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	duration := time.Since(commandStartTime)

	// End the root span started in initTelemetry.
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))
	span.End()

	telemetry.Logger.InfoContext(ctx, "command.finish",
		"command", cmd.CommandPath(),
		"duration_ms", duration.Milliseconds(),
	)

	telemetry.RecordCommandMetrics(ctx, cmd.CommandPath(), duration, true)

	if telemetryShutdown != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return telemetryShutdown(shutdownCtx)
	}
	return nil
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".allbctl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
