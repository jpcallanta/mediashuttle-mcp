# Signiant Media Shuttle MCP

MCP server for managing Signiant Media Shuttle via an AI agent.

## Architecture

```
AI Agent / MCP Client  <--stdio-->  mediashuttle-mcp (MCP Server)
                                       |
                                  internal/server
                                       |  (13 MCP Tools)
                                  internal/client
                                       |  (HTTP REST)
                           https://api.mediashuttle.com/v1
```

## Usage

```sh
# Set your API key
export MS_API_KEY=your_api_key_here

# Run the MCP server (stdio transport)
mediashuttle-mcp

# Or pass the key as a flag
mediashuttle-mcp --key your_api_key_here

# Demo mode (exercises read-only API calls without an MCP client)
mediashuttle-mcp --demo
```

## MCP Tools

| Category       | Tool                   | Description                     |
|----------------|------------------------|---------------------------------|
| **Portals**    | `list_portals`         | List all portals                |
|                | `create_portal`        | Create a portal                 |
|                | `update_portal`        | Update a portal                 |
| **Portal
Users**    | `list_portal_users`    | List users in a portal          |
|                | `get_portal_user`      | Get a portal user by email      |
|                | `add_portal_user`      | Add a user to a portal          |
|                | `update_portal_user`   | Update a portal user            |
|                | `remove_portal_user`   | Remove a user from a portal     |
| **Portal
Storage** | `list_portal_storage`  | List storage assigned to portal |
|                | `assign_portal_storage`| Assign storage to a portal      |
| **Storage**    | `list_storage`         | List all storage nodes          |
|                | `get_storage`          | Get storage details             |
| **Transfers**  | `list_transfers`       | List active transfers           |

## Client Setup

### Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`
(macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "mediashuttle": {
      "command": "/path/to/mediashuttle-mcp",
      "args": [],
      "env": {
        "MS_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

### Claude Code

```sh
# Add as a local stdio server
claude mcp add --transport stdio mediashuttle -- \
  /path/to/mediashuttle-mcp

# Or with an env var
claude mcp add --env MS_API_KEY=your_api_key \
  --transport stdio mediashuttle -- /path/to/mediashuttle-mcp
```

Or create `.mcp.json` in your project root:

```json
{
  "mcpServers": {
    "mediashuttle": {
      "command": "/path/to/mediashuttle-mcp",
      "args": [],
      "env": {
        "MS_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

### OpenCode

Add to your `opencode.json`:

```json
{
  "mcp": {
    "mediashuttle": {
      "type": "local",
      "command": ["/path/to/mediashuttle-mcp"],
      "enabled": true,
      "environment": {
        "MS_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

### VS Code

Edit `.vscode/mcp.json` in your workspace (or user-level via
**MCP: Open User Configuration**):

```json
{
  "servers": {
    "mediashuttle": {
      "type": "stdio",
      "command": "/path/to/mediashuttle-mcp",
      "args": [],
      "env": {
        "MS_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

You can also install via the Extensions view (`@mcp` search) if
published, or use **MCP: Add Server** from the Command Palette.

### Gemini / Other MCP Clients

Any MCP-compatible client accepts the same stdio configuration.
Point the client's MCP server command at the `mediashuttle-mcp`
binary with `MS_API_KEY` set in the environment.

## Development

```sh
# Build
make

# Test
make test

# Install (auto-detects /usr/local/bin or ~/go/bin)
make install

# Override install prefix
make install PREFIX=/opt/local

# For package builds
make install DESTDIR=/tmp/pkg
```

### Install locations

`make install` copies the binary to:

1. `/usr/local/bin` — if you have write permission (e.g. `sudo make
   install` or a writable `/usr/local/bin`)
2. `~/go/bin` — fallback when `/usr/local/bin` is not writable

Override with `PREFIX=<dir>` or `DESTDIR=<dir>` for staged
installations.
