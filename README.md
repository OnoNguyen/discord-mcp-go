# discord-mcp-go

A Go-based [MCP](https://modelcontextprotocol.io/) server that exposes Discord operations as tools. Use it with Claude Code (or any MCP client) to read channels, threads, and messages — enabling TLDR summaries, channel digests, and more.

## Features

31 tools across 9 categories:

| Category | Tools |
|----------|-------|
| **Server** | `get_server_info` |
| **Messages** | `read_messages`, `send_message`, `edit_message`, `delete_message` |
| **Reactions** | `add_reaction`, `remove_reaction` |
| **Channels** | `list_channels`, `find_channel`, `create_text_channel`, `delete_channel` |
| **Threads** | `list_threads`, `read_thread_messages` |
| **Categories** | `list_categories`, `find_category`, `create_category`, `delete_category`, `list_channels_in_category` |
| **Users** | `get_user_id_by_name`, `send_private_message`, `read_private_messages` |
| **Roles** | `list_roles`, `create_role`, `edit_role`, `delete_role`, `assign_role`, `remove_role` |
| **Webhooks** | `list_webhooks`, `create_webhook`, `delete_webhook`, `send_webhook_message` |

## Prerequisites

1. **Go 1.23+** (for building from source)
2. **Discord Bot Token** — create one at [discord.com/developers](https://discord.com/developers/applications)
3. **Bot Permissions** — the bot needs permissions for the operations you want to use (Read Messages, Send Messages, Manage Roles, etc.)

### Bot Setup

1. Create an application at [Discord Developer Portal](https://discord.com/developers/applications)
2. Go to **Bot** tab, click **Reset Token**, and copy it
3. Enable **Message Content Intent** under Privileged Gateway Intents
4. Go to **OAuth2 > URL Generator**, select `bot` scope with needed permissions
5. Use the generated URL to invite the bot to your server

## Installation

### Build from source

```bash
git clone https://github.com/OnoNguyen/discord-mcp-go.git
cd discord-mcp-go
go build -o discord-mcp-go .
```

### Docker

```bash
docker build -t discord-mcp-go .
```

## Configuration

Set environment variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `DISCORD_TOKEN` | Yes | Discord bot token |
| `DISCORD_GUILD_ID` | No | Default server/guild ID (avoids passing it to every tool) |

## Usage with Claude Code

### Binary

```bash
claude mcp add discord-mcp -- \
  env DISCORD_TOKEN=your-token DISCORD_GUILD_ID=your-guild-id \
  /path/to/discord-mcp-go
```

### Docker

```bash
claude mcp add discord-mcp -- \
  docker run --rm -i \
  -e DISCORD_TOKEN=your-token \
  -e DISCORD_GUILD_ID=your-guild-id \
  discord-mcp-go
```

### Then in Claude Code

```
> list all channels in my Discord server

> read the last 50 messages in #general and give me a TLDR

> summarize the active threads in #engineering
```

## Development

```bash
# Build
go build -o discord-mcp-go .

# Run directly (for testing)
DISCORD_TOKEN=xxx DISCORD_GUILD_ID=xxx ./discord-mcp-go
```

## License

MIT
