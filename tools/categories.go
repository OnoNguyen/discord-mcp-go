package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListCategoriesInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
}

type FindCategoryInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	Name    string `json:"name" jsonschema:"Category name to search for"`
}

type CreateCategoryInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	Name    string `json:"name" jsonschema:"Name for the new category"`
}

type DeleteCategoryInput struct {
	CategoryID string `json:"categoryId" jsonschema:"ID of the category to delete"`
}

type ListChannelsInCategoryInput struct {
	CategoryID string `json:"categoryId" jsonschema:"Category ID to list channels from"`
	GuildID    string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
}

func RegisterCategoryTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all categories in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListCategoriesInput) (*mcp.CallToolResult, any, error) {
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
			if ch.Type == discordgo.ChannelTypeGuildCategory {
				sb.WriteString(fmt.Sprintf("📁 %s (ID: %s)\n", ch.Name, ch.ID))
			}
		}

		if sb.Len() == 0 {
			return textResult("No categories found."), nil, nil
		}
		return textResult(sb.String()), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "find_category",
		Description: "Find a category by name in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input FindCategoryInput) (*mcp.CallToolResult, any, error) {
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
			if ch.Type == discordgo.ChannelTypeGuildCategory && strings.Contains(strings.ToLower(ch.Name), searchName) {
				matches = append(matches, fmt.Sprintf("📁 %s (ID: %s)", ch.Name, ch.ID))
			}
		}

		if len(matches) == 0 {
			return textResult(fmt.Sprintf("No categories matching '%s' found.", input.Name)), nil, nil
		}
		return textResult(strings.Join(matches, "\n")), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_category",
		Description: "Create a new category in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateCategoryInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		ch, err := dc.Session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: input.Name,
			Type: discordgo.ChannelTypeGuildCategory,
		})
		if err != nil {
			return errResult(fmt.Errorf("creating category: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Category '%s' created (ID: %s)", ch.Name, ch.ID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_category",
		Description: "Delete a category from a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteCategoryInput) (*mcp.CallToolResult, any, error) {
		ch, err := dc.Session.ChannelDelete(input.CategoryID)
		if err != nil {
			return errResult(fmt.Errorf("deleting category: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Category '%s' deleted", ch.Name)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_channels_in_category",
		Description: "List all channels within a specific category",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListChannelsInCategoryInput) (*mcp.CallToolResult, any, error) {
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
			if ch.ParentID == input.CategoryID {
				sb.WriteString(fmt.Sprintf("%s %s (ID: %s, Type: %s)\n", channelTypeIcon(ch.Type), ch.Name, ch.ID, channelTypeName(ch.Type)))
			}
		}

		if sb.Len() == 0 {
			return textResult("No channels found in this category."), nil, nil
		}
		return textResult(sb.String()), nil, nil
	})
}
