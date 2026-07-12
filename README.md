# Signiant Media Shuttle MCP

MCP server for managing Signiant Media Shuttle via an AI agent.

## Architecture

```
                    ┌── stdio ──→  mediashuttle-mcp (stdio)
                    │                |
AI Agent /          │           internal/server
MCP Client ─────────┤                |  (13 MCP Tools)
                    │            internal/client
                    │                |  (HTTP REST)
                    └── HTTP ──→  mediashuttle-mcp serve
                         |     (Streamable HTTP, port 8080)
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

# Run the MCP server (stdio transport — single client)
mediashuttle-mcp

# Or pass the key as a flag
mediashuttle-mcp --key your_api_key_here

# Demo mode (exercises read-only API calls without an MCP client)
mediashuttle-mcp --demo
```

### HTTP Server (Streamable HTTP Transport)

The `serve` subcommand starts an HTTP server using the
[MCP Streamable HTTP](https://modelcontextprotocol.io/specification/2025-03-26/basic/transports#streamable-http)
transport, supporting **multiple concurrent clients** over a single
process.

```sh
# Start on default port :8080
mediashuttle-mcp serve

# Custom address
mediashuttle-mcp serve --addr :9090

# With API key flag
mediashuttle-mcp --key your_key_here serve --addr :8080
```

The HTTP server exposes a single endpoint (`/mcp` by default):

| Method   | Purpose                                |
|----------|----------------------------------------|
| `POST`   | JSON-RPC requests (initialize, tools)  |
| `GET`    | SSE stream for notifications           |
| `DELETE` | Session cleanup                        |

Clients receive a `Mcp-Session-Id` header on initialize and must
include it on subsequent requests.

### Docker

Build and run the HTTP server via Docker:

```sh
# Build the image
make docker-image

# Or build manually
docker build -t mediashuttle-mcp .
```

Run with docker-compose (recommended):

```sh
export MS_API_KEY=your_api_key_here
make docker-up
```

Or run directly:

```sh
docker run -d --name mediashuttle -p 8080:8080 \
  -e MS_API_KEY=your_api_key_here \
  mediashuttle-mcp
```

See `Dockerfile` and `docker-compose.yml` in the project root.

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

### Streamable HTTP Setup

These examples assume the server is running on `http://localhost:8080`.
Start it with:

```sh
mediashuttle-mcp serve
```

### Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`
(macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

**stdio** (single process):
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

**Streamable HTTP** (shared server):
```json
{
  "mcpServers": {
    "mediashuttle": {
      "type": "sse",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### Claude Code

**stdio** (single process):
```sh
claude mcp add --transport stdio mediashuttle -- \
  /path/to/mediashuttle-mcp

claude mcp add --env MS_API_KEY=your_api_key \
  --transport stdio mediashuttle -- /path/to/mediashuttle-mcp
```

**Streamable HTTP** (shared server):
```sh
claude mcp add --transport sse mediashuttle -- \
  http://localhost:8080/mcp
```

Or create `.mcp.json` in your project root:

**stdio:**
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

**Streamable HTTP:**
```json
{
  "mcpServers": {
    "mediashuttle": {
      "type": "sse",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### OpenCode

**stdio** (single process) — add to `opencode.json`:
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

**Streamable HTTP** (shared server):
```json
{
  "mcp": {
    "mediashuttle": {
      "type": "sse",
      "url": "http://localhost:8080/mcp",
      "enabled": true
    }
  }
}
```

### VS Code

Edit `.vscode/mcp.json` in your workspace (or user-level via
**MCP: Open User Configuration**):

**stdio:**
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

**Streamable HTTP:**
```json
{
  "servers": {
    "mediashuttle": {
      "type": "sse",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

You can also install via the Extensions view (`@mcp` search) if
published, or use **MCP: Add Server** from the Command Palette.

### Gemini / Other MCP Clients

**stdio** — point the client's MCP server command at the
`mediashuttle-mcp` binary with `MS_API_KEY` in the environment.

**Streamable HTTP** — configure the client with the SSE/HTTP
transport pointing at `http://<host>:8080/mcp`. The MCP
Streamable HTTP transport is widely supported; consult your
client's documentation for the exact configuration key.

## Development

```sh
# Build
make

# Test
make test

# Install (auto-detects /usr/local/bin or ~/go/bin)
make install

# Docker image
make docker-image

# Docker compose (build + run)
make docker-up

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
