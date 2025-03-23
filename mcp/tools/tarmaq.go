package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/mazrean/mcp-tarmaq/tarmaq"
)

var _ Tool = &TarmaqTool{}

type TarmaqTool struct {
	executer *tarmaq.Tarmaq
}

func NewTarmaqTool(executer *tarmaq.Tarmaq) *TarmaqTool {
	return &TarmaqTool{
		executer: executer,
	}
}

func (h *TarmaqTool) Tool() mcp.Tool {
	return mcp.NewTool("impact_analysis",
		mcp.WithDescription("Suggest files that are likely to change at the same time in the changelog"),
		mcp.WithArray("files",
			mcp.Required(),
			mcp.Description("already modified files"),
		),
		mcp.WithNumber("min_confidence",
			mcp.DefaultNumber(0.7),
			mcp.Description("minimum confidence. confidence represents the probability of change(0.0 ~ 1.0)."),
		),
		mcp.WithNumber("min_support",
			mcp.DefaultNumber(3),
			mcp.Description("minimum support. support represents the number of occurrences of the change."),
		),
	)
}

type TarmaqResponse struct {
	Path       string  `json:"file_path"`
	Confidence float64 `json:"confidence"`
	Support    uint64  `json:"support"`
}

func (h *TarmaqTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	iFiles, ok := request.Params.Arguments["files"].([]any)
	if !ok {
		slog.Error("invalid files",
			slog.String("files", fmt.Sprintf("%v", request.Params.Arguments["files"])),
		)
		return nil, fmt.Errorf("invalid files: %v", request.Params.Arguments["files"])
	}

	minConfidence, ok := request.Params.Arguments["min_confidence"].(float64)
	if !ok {
		slog.Error("invalid min_confidence",
			slog.String("min_confidence", fmt.Sprintf("%v", request.Params.Arguments["min_confidence"])),
		)
		return nil, fmt.Errorf("invalid min_confidence: %v", request.Params.Arguments["min_confidence"])
	}

	minSupport, ok := request.Params.Arguments["min_support"].(float64)
	if !ok {
		slog.Error("invalid min_support",
			slog.String("min_support", fmt.Sprintf("%v", request.Params.Arguments["min_support"])),
		)
		return nil, fmt.Errorf("invalid min_support: %v", request.Params.Arguments["min_support"])
	}

	tarmaqFiles := make([]tarmaq.FilePath, 0, len(iFiles))
	for _, iFile := range iFiles {
		file, ok := iFile.(string)
		if !ok {
			slog.Warn("invalid file",
				slog.String("file", fmt.Sprintf("%v", iFile)),
			)
			continue
		}
		tarmaqFiles = append(tarmaqFiles, tarmaq.FilePath(file))
	}
	result, err := h.executer.Execute(tarmaqFiles, minConfidence, uint64(minSupport))
	if err != nil {
		slog.Error("execute tarmaq",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("execute tarmaq: %w", err)
	}

	res := make([]*TarmaqResponse, 0, len(result))
	for _, rule := range result {
		res = append(res, &TarmaqResponse{
			Path:       filepath.FromSlash(string(rule.Path)),
			Confidence: rule.Confidence,
			Support:    rule.Support,
		})
	}

	response, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		slog.Error("marshal response",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(response)), nil
}
