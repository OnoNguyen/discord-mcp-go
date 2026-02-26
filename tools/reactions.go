package tools

import (
	"context"
	"fmt"

	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddReactionInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID containing the message"`
	MessageID string `json:"messageId" jsonschema:"ID of the message to react to"`
	Emoji     string `json:"emoji" jsonschema:"Emoji to react with (e.g. 👍 or custom emoji name:id)"`
}

type RemoveReactionInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID containing the message"`
	MessageID string `json:"messageId" jsonschema:"ID of the message"`
	Emoji     string `json:"emoji" jsonschema:"Emoji to remove (e.g. 👍 or custom emoji name:id)"`
}

func RegisterReactionTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_reaction",
		Description: "Add a reaction emoji to a message",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AddReactionInput) (*mcp.CallToolResult, any, error) {
		err := dc.Session.MessageReactionAdd(input.ChannelID, input.MessageID, input.Emoji)
		if err != nil {
			return errResult(fmt.Errorf("adding reaction: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Reaction %s added to message %s", input.Emoji, input.MessageID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_reaction",
		Description: "Remove the bot's reaction from a message",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input RemoveReactionInput) (*mcp.CallToolResult, any, error) {
		err := dc.Session.MessageReactionRemove(input.ChannelID, input.MessageID, input.Emoji, "@me")
		if err != nil {
			return errResult(fmt.Errorf("removing reaction: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Reaction %s removed from message %s", input.Emoji, input.MessageID)), nil, nil
	})
}
