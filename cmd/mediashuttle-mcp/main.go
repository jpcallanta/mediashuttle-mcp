package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"mediashuttle-mcp/internal/client"
	msserver "mediashuttle-mcp/internal/server"
)

var apiKey string
var demo bool

func init() {
	root.Flags().StringVar(
		&apiKey, "key", "",
		"Media Shuttle API key"+
			" (overrides MS_API_KEY env var)",
	)
	root.Flags().BoolVar(
		&demo, "demo", false,
		"Run capability demo instead of MCP server",
	)
}

var root = &cobra.Command{
	Use:   "mediashuttle-mcp",
	Short: "Media Shuttle MCP server",
	Long: "A Model Context Protocol server" +
		" for managing Signiant Media Shuttle" +
		" users, portals, storage, and transfers.",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if demo {
			c := client.NewClient(apiKey)
			RunDemo(c)
			return nil
		}

		s := msserver.NewMCPServer(apiKey)
		return server.ServeStdio(s)
	},
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
