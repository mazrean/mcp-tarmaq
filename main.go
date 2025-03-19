package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

var (
	version  = "dev"
	revision = "none"
)

// CLI represents command line options and configuration file values
var CLI struct {
	Version  kong.VersionFlag `kong:"short='v',help='Show version and exit.'"`
	LogLevel string           `kong:"short='l',default='info',enum='debug,info,warn,error',help='Log level',env='MCP_TARMAQ_LOG_LEVEL'"`
}

// loadConfig loads and parses configuration from command line arguments
func loadConfig() (*kong.Context, error) {
	// Parse command line arguments
	parser := kong.Must(&CLI,
		kong.Name("mcp-tarmaq"),
		kong.Description("MCP server for impact analysis"),
		kong.Vars{"version": fmt.Sprintf("%s (%s)", version, revision)},
		kong.UsageOnError(),
	)
	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	return ctx, nil
}

func main() {
	_, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var level slog.Level
	switch CLI.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))
}
