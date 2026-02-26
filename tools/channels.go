package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListChannelsInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
}

type FindChannelInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	Name    string `json:"name" jsonschema:"Channel name to search for"`
}

type CreateTextChannelInput struct {
	GuildID    string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	Name       string `json:"name" jsonschema:"Name for the new channel"`
	CategoryID string `json:"categoryId,omitempty" jsonschema:"Category ID to place the channel in"`
}

type DeleteChannelInput struct {
	ChannelID string `json:"channelId" jsonschema:"ID of the channel to delete"`
}

func RegisterChannelTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_channels",
		Description: "List all channels in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListChannelsInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		channels, err := dc.Session.GuildChannels(guildID)
		if err != nil {
			return errResult(fmt.Errorf("fetching channels: %w", err)), nil, nil
		}

		var sb strings.Builder
		for _, ch := range channels {
			chType := channelTypeName(ch.Type)
			sb.WriteString(fmt.Sprintf("%s %s (ID: %s, Type: %s)\n", channelTypeIcon(ch.Type), ch.Name, ch.ID, chType))
		}

		if sb.Len() == 0 {
			return textResult("No channels found."), nil, nil
		}
		return textResult(sb.String()), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "find_channel",
		Description: "Find a channel by name in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input FindChannelInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		channels, err := dc.Session.GuildChannels(guildID)
		if err != nil {
			return errResult(fmt.Errorf("fetching channels: %w", err)), nil, nil
		}

		searchName := strings.ToLower(input.Name)
		var matches []string
		for _, ch := range channels {
			if strings.Contains(strings.ToLower(ch.Name), searchName) {
				matches = append(matches, fmt.Sprintf("%s %s (ID: %s, Type: %s)", channelTypeIcon(ch.Type), ch.Name, ch.ID, channelTypeName(ch.Type)))
			}
		}

		if len(matches) == 0 {
			return textResult(fmt.Sprintf("No channels matching '%s' found.", input.Name)), nil, nil
		}
		return textResult(strings.Join(matches, "\n")), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_text_channel",
		Description: "Create a new text channel in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateTextChannelInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		data := discordgo.GuildChannelCreateData{
			Name: input.Name,
			Type: discordgo.ChannelTypeGuildText,
		}
		if input.CategoryID != "" {
			data.ParentID = input.CategoryID
		}

		ch, err := dc.Session.GuildChannelCreateComplex(guildID, data)
		if err != nil {
			return errResult(fmt.Errorf("creating channel: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Channel #%s created (ID: %s)", ch.Name, ch.ID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_channel",
		Description: "Delete a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteChannelInput) (*mcp.CallToolResult, any, error) {
		ch, err := dc.Session.ChannelDelete(input.ChannelID)
		if err != nil {
			return errResult(fmt.Errorf("deleting channel: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Channel #%s deleted", ch.Name)), nil, nil
	})
}

func channelTypeName(t discordgo.ChannelType) string {
	switch t {
	case discordgo.ChannelTypeGuildText:
		return "text"
	case discordgo.ChannelTypeGuildVoice:
		return "voice"
	case discordgo.ChannelTypeGuildCategory:
		return "category"
	case discordgo.ChannelTypeGuildNews:
		return "news"
	case discordgo.ChannelTypeGuildStageVoice:
		return "stage"
	case discordgo.ChannelTypeGuildForum:
		return "forum"
	default:
		return fmt.Sprintf("type-%d", t)
	}
}

func channelTypeIcon(t discordgo.ChannelType) string {
	switch t {
	case discordgo.ChannelTypeGuildText:
		return "#"
	case discordgo.ChannelTypeGuildVoice:
		return "🔊"
	case discordgo.ChannelTypeGuildCategory:
		return "📁"
	case discordgo.ChannelTypeGuildNews:
		return "📢"
	case discordgo.ChannelTypeGuildStageVoice:
		return "🎙️"
	case discordgo.ChannelTypeGuildForum:
		return "💬"
	default:
		return "•"
	}
}
