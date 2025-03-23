package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/mazrean/mcp-tarmaq/mcp"
	"github.com/mazrean/mcp-tarmaq/mcp/tools"
	"github.com/mazrean/mcp-tarmaq/tarmaq"
)

var (
	version  = "dev"
	revision = "none"
)

// CLI represents command line options and configuration file values
var CLI struct {
	Version        kong.VersionFlag `kong:"short='v',help='Show version and exit.'"`
	LogLevel       string           `kong:"short='l',default='info',enum='debug,info,warn,error',help='Log level',env='MCP_TARMAQ_LOG_LEVEL'"`
	RepositoryPath string           `kong:"short='r',help='Path to the repository',env='MCP_TARMAQ_REPOSITORY_PATH'"`
	CommitLimit    int              `kong:"default='200',help='Limit of commits to analyze',env='MCP_TARMAQ_COMMIT_LIMIT'"`
	MaxChangedFile int              `kong:"default='30',help='Limit of changed files in a commit',env='MCP_TARMAQ_MAX_CHANGED_FILE'"`
}

// loadConfig loads and parses configuration from command line arguments
func loadConfig() (*kong.Context, error) {
	// Parse command line arguments
	parser := kong.Must(&CLI,
		kong.Name("mcp-tarmaq"),
		kong.Description("MCP server for impact analysis.Suggest files that are likely to change at the same time in the changelog."),
		kong.Vars{"version": fmt.Sprintf("%s (%s)", version, revision)},
		kong.UsageOnError(),
	)
	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	return ctx, nil
}

func createTarmaq() (*tarmaq.Tarmaq, error) {
	repo, err := tarmaq.NewGitRepository(CLI.RepositoryPath, CLI.CommitLimit)
	if err != nil {
		return nil, fmt.Errorf("create git repository: %w", err)
	}

	tarmaq := tarmaq.NewTarmaq(repo, []tarmaq.TxFilter{
		tarmaq.NewMaxSizeTxFilter(CLI.MaxChangedFile),
		tarmaq.NewTarmaqTxFilter(),
	}, tarmaq.NewAssociationRuleExtractor())

	return tarmaq, nil
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

	executer, err := createTarmaq()
	if err != nil {
		slog.Error("failed to create tarmaq",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	server := mcp.NewServer(version, tools.NewTarmaqTool(executer))
	if err := server.Start(); err != nil {
		slog.Error("failed to run server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
