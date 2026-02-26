package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListRolesInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
}

type CreateRoleInput struct {
	GuildID     string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	Name        string `json:"name" jsonschema:"Name for the new role"`
	Color       int    `json:"color,omitempty" jsonschema:"Role color as integer (e.g. 0xFF0000 for red)"`
	Permissions int64  `json:"permissions,omitempty" jsonschema:"Permission bitfield for the role"`
}

type EditRoleInput struct {
	GuildID     string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	RoleID      string `json:"roleId" jsonschema:"ID of the role to edit"`
	Name        string `json:"name,omitempty" jsonschema:"New name for the role"`
	Color       int    `json:"color,omitempty" jsonschema:"New color as integer"`
	Permissions int64  `json:"permissions,omitempty" jsonschema:"New permission bitfield"`
}

type DeleteRoleInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	RoleID  string `json:"roleId" jsonschema:"ID of the role to delete"`
}

type AssignRoleInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	UserID  string `json:"userId" jsonschema:"User ID to assign the role to"`
	RoleID  string `json:"roleId" jsonschema:"Role ID to assign"`
}

type RemoveRoleInput struct {
	GuildID string `json:"guildId,omitempty" jsonschema:"Server/guild ID (uses default if not provided)"`
	UserID  string `json:"userId" jsonschema:"User ID to remove the role from"`
	RoleID  string `json:"roleId" jsonschema:"Role ID to remove"`
}

func RegisterRoleTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_roles",
		Description: "List all roles in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListRolesInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		roles, err := dc.Session.GuildRoles(guildID)
		if err != nil {
			return errResult(fmt.Errorf("fetching roles: %w", err)), nil, nil
		}

		var sb strings.Builder
		for _, r := range roles {
			sb.WriteString(fmt.Sprintf("@%s (ID: %s, Color: #%06x, Members: managed=%v)\n", r.Name, r.ID, r.Color, r.Managed))
		}
		return textResult(sb.String()), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_role",
		Description: "Create a new role in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateRoleInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		params := &discordgo.RoleParams{
			Name: input.Name,
		}
		if input.Color != 0 {
			params.Color = &input.Color
		}
		if input.Permissions != 0 {
			params.Permissions = &input.Permissions
		}

		role, err := dc.Session.GuildRoleCreate(guildID, params)
		if err != nil {
			return errResult(fmt.Errorf("creating role: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Role @%s created (ID: %s)", role.Name, role.ID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_role",
		Description: "Edit an existing role in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input EditRoleInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		params := &discordgo.RoleParams{}
		if input.Name != "" {
			params.Name = input.Name
		}
		if input.Color != 0 {
			params.Color = &input.Color
		}
		if input.Permissions != 0 {
			params.Permissions = &input.Permissions
		}

		role, err := dc.Session.GuildRoleEdit(guildID, input.RoleID, params)
		if err != nil {
			return errResult(fmt.Errorf("editing role: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Role @%s updated (ID: %s)", role.Name, role.ID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_role",
		Description: "Delete a role from a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteRoleInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		err = dc.Session.GuildRoleDelete(guildID, input.RoleID)
		if err != nil {
			return errResult(fmt.Errorf("deleting role: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Role %s deleted", input.RoleID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "assign_role",
		Description: "Assign a role to a user in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AssignRoleInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		err = dc.Session.GuildMemberRoleAdd(guildID, input.UserID, input.RoleID)
		if err != nil {
			return errResult(fmt.Errorf("assigning role: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Role %s assigned to user %s", input.RoleID, input.UserID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_role",
		Description: "Remove a role from a user in a Discord server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input RemoveRoleInput) (*mcp.CallToolResult, any, error) {
		guildID, err := dc.ResolveGuildID(input.GuildID)
		if err != nil {
			return errResult(err), nil, nil
		}

		err = dc.Session.GuildMemberRoleRemove(guildID, input.UserID, input.RoleID)
		if err != nil {
			return errResult(fmt.Errorf("removing role: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Role %s removed from user %s", input.RoleID, input.UserID)), nil, nil
	})
}
