package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	msserver "mediashuttle-mcp/internal/server"
)

var serveAddr string

func init() {
	serveCmd.Flags().StringVar(
		&serveAddr, "addr", ":8080",
		"Listen address for HTTP server",
	)
	root.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start MCP server over HTTP",
	Long: "Starts an HTTP server that serves the MCP tools" +
		" over the Streamable HTTP transport," +
		" allowing multiple clients to connect.",
	RunE: func(
		cmd *cobra.Command, args []string,
	) error {
		if apiKey == "" {
			apiKey = os.Getenv("MS_API_KEY")
		}

		if apiKey == "" {
			return fmt.Errorf(
				"API key required:" +
					" use --key flag" +
					" or set MS_API_KEY env var",
			)
		}

		s := msserver.NewMCPServer(apiKey)
		httpSrv := server.NewStreamableHTTPServer(s)
		return httpSrv.Start(serveAddr)
	},
}
