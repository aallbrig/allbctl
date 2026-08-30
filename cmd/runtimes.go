package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/aallbrig/allbctl/pkg/telemetry"
)

var RuntimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Display detected runtimes and their versions",
	Long: `Display detected programming language runtimes and their versions.

This is the same output shown in the 'Runtimes:' section of 'allbctl status'.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, span := otel.Tracer("github.com/aallbrig/allbctl").Start(ctx, "status.runtimes")
		defer span.End()
		PrintRuntimes(ctx)
	},
}

// PrintRuntimes outputs the runtimes in inline format (same as status command)
func PrintRuntimes(ctx context.Context) {
	runtimesInline := detectRuntimesInline()
	if runtimesInline != "" {
		fmt.Println(runtimesInline)
	} else {
		fmt.Println("No runtimes detected")
	}

	count := 0
	if runtimesInline != "" {
		count = len(strings.Split(runtimesInline, ","))
	}
	telemetry.Logger.InfoContext(ctx, "status.runtimes",
		"runtimes", runtimesInline,
		"count", count,
	)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("runtime_count", count))
}
