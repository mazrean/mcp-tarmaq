package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/mazrean/mcp-tarmaq/mcp/tools"
)

type Server struct {
	server *server.MCPServer
}

func NewServer(
	version string,
	tarmaqTool tools.Tool,
) *Server {
	s := server.NewMCPServer(
		"tarmaq", // Name
		version,  // Version
		server.WithLogging(),
	)

	s.AddTool(tarmaqTool.Tool(), tarmaqTool.Handle)

	return &Server{
		server: s,
	}
}

func (s *Server) Start() error {
	return server.ServeStdio(s.server)
}
