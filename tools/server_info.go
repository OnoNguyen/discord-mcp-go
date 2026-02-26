package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetServerInfoInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
}

func RegisterServerInfoTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_server_info",
		Description: "Get Discord server name, member count, channels, and roles overview",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetServerInfoInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		guild, err := dc.Session.Guild(guildID)
		if err != nil {
			return errResult(fmt.Errorf("fetching guild: %w", err)), nil, nil
		}

		channels, err := dc.Session.GuildChannels(guildID)
		if err != nil {
			return errResult(fmt.Errorf("fetching channels: %w", err)), nil, nil
		}

		roles, err := dc.Session.GuildRoles(guildID)
		if err != nil {
			return errResult(fmt.Errorf("fetching roles: %w", err)), nil, nil
		}

		var textChs, voiceChs, categories []string
		for _, ch := range channels {
			switch ch.Type {
			case 0: // text
				textChs = append(textChs, fmt.Sprintf("#%s (ID: %s)", ch.Name, ch.ID))
			case 2: // voice
				voiceChs = append(voiceChs, fmt.Sprintf("🔊 %s (ID: %s)", ch.Name, ch.ID))
			case 4: // category
				categories = append(categories, fmt.Sprintf("📁 %s (ID: %s)", ch.Name, ch.ID))
			}
		}

		var roleNames []string
		for _, r := range roles {
			roleNames = append(roleNames, fmt.Sprintf("@%s (ID: %s)", r.Name, r.ID))
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("**Server:** %s\n", guild.Name))
		sb.WriteString(fmt.Sprintf("**ID:** %s\n", guild.ID))
		sb.WriteString(fmt.Sprintf("**Member Count:** %d\n", guild.MemberCount))
		sb.WriteString(fmt.Sprintf("**Owner ID:** %s\n\n", guild.OwnerID))

		sb.WriteString(fmt.Sprintf("**Categories (%d):**\n%s\n\n", len(categories), strings.Join(categories, "\n")))
		sb.WriteString(fmt.Sprintf("**Text Channels (%d):**\n%s\n\n", len(textChs), strings.Join(textChs, "\n")))
		sb.WriteString(fmt.Sprintf("**Voice Channels (%d):**\n%s\n\n", len(voiceChs), strings.Join(voiceChs, "\n")))
		sb.WriteString(fmt.Sprintf("**Roles (%d):**\n%s", len(roleNames), strings.Join(roleNames, "\n")))

		return textResult(sb.String()), nil, nil
	})
}
