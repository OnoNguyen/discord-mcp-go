package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetUserIDByNameInput struct {
	GuildID  string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	Username string `json:"username" jsonschema:"Username to search for"`
}

type SendPrivateMessageInput struct {
	UserID  string `json:"userId" jsonschema:"User ID to send the DM to"`
	Content string `json:"content" jsonschema:"Message content"`
}

type ReadPrivateMessagesInput struct {
	UserID string `json:"userId" jsonschema:"User ID to read DMs from"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Number of messages to fetch (default 50, max 100)"`
}

func RegisterUserTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_user_id_by_name",
		Description: "Look up a Discord user's ID by their username in a server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetUserIDByNameInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		// Search guild members
		members, err := dc.Session.GuildMembersSearch(guildID, input.Username, 10)
		if err != nil {
			return errResult(fmt.Errorf("searching members: %w", err)), nil, nil
		}

		if len(members) == 0 {
			return textResult(fmt.Sprintf("No users matching '%s' found.", input.Username)), nil, nil
		}

		var sb strings.Builder
		for _, m := range members {
			nick := ""
			if m.Nick != "" {
				nick = fmt.Sprintf(" (nick: %s)", m.Nick)
			}
			sb.WriteString(fmt.Sprintf("%s#%s%s — ID: %s\n", m.User.Username, m.User.Discriminator, nick, m.User.ID))
		}
		return textResult(sb.String()), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_private_message",
		Description: "Send a direct message to a Discord user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SendPrivateMessageInput) (*mcp.CallToolResult, any, error) {
		channel, err := dc.Session.UserChannelCreate(input.UserID)
		if err != nil {
			return errResult(fmt.Errorf("creating DM channel: %w", err)), nil, nil
		}

		msg, err := dc.Session.ChannelMessageSend(channel.ID, input.Content)
		if err != nil {
			return errResult(fmt.Errorf("sending DM: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("DM sent (Message ID: %s)", msg.ID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_private_messages",
		Description: "Read direct message history with a Discord user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ReadPrivateMessagesInput) (*mcp.CallToolResult, any, error) {
		channel, err := dc.Session.UserChannelCreate(input.UserID)
		if err != nil {
			return errResult(fmt.Errorf("creating DM channel: %w", err)), nil, nil
		}

		limit := input.Limit
		if limit <= 0 || limit > 100 {
			limit = 50
		}

		messages, err := dc.Session.ChannelMessages(channel.ID, limit, "", "", "")
		if err != nil {
			return errResult(fmt.Errorf("fetching DMs: %w", err)), nil, nil
		}

		return textResult(formatMessages(messages)), nil, nil
	})
}
