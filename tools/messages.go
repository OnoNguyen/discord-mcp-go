package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ReadMessagesInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID to read messages from"`
	Limit     int    `json:"limit,omitempty" jsonschema:"Number of messages to fetch (default 100, max 100)"`
}

type SendMessageInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID to send the message to"`
	Content   string `json:"content" jsonschema:"Message content"`
}

type EditMessageInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID containing the message"`
	MessageID string `json:"messageId" jsonschema:"ID of the message to edit"`
	Content   string `json:"content" jsonschema:"New message content"`
}

type DeleteMessageInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID containing the message"`
	MessageID string `json:"messageId" jsonschema:"ID of the message to delete"`
}

func RegisterMessageTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_messages",
		Description: "Read recent messages from a Discord channel. Returns messages in chronological order with author, timestamp, and content.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ReadMessagesInput) (*mcp.CallToolResult, any, error) {
		limit := input.Limit
		if limit <= 0 || limit > 100 {
			limit = 100
		}

		messages, err := dc.Session.ChannelMessages(input.ChannelID, limit, "", "", "")
		if err != nil {
			return errResult(fmt.Errorf("fetching messages: %w", err)), nil, nil
		}

		return textResult(formatMessages(messages)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_message",
		Description: "Send a message to a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SendMessageInput) (*mcp.CallToolResult, any, error) {
		msg, err := dc.Session.ChannelMessageSend(input.ChannelID, input.Content)
		if err != nil {
			return errResult(fmt.Errorf("sending message: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Message sent (ID: %s)", msg.ID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_message",
		Description: "Edit an existing message in a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input EditMessageInput) (*mcp.CallToolResult, any, error) {
		_, err := dc.Session.ChannelMessageEdit(input.ChannelID, input.MessageID, input.Content)
		if err != nil {
			return errResult(fmt.Errorf("editing message: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Message %s edited successfully", input.MessageID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_message",
		Description: "Delete a message from a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteMessageInput) (*mcp.CallToolResult, any, error) {
		err := dc.Session.ChannelMessageDelete(input.ChannelID, input.MessageID)
		if err != nil {
			return errResult(fmt.Errorf("deleting message: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Message %s deleted successfully", input.MessageID)), nil, nil
	})
}

func formatMessages(messages []*discordgo.Message) string {
	if len(messages) == 0 {
		return "No messages found."
	}

	// Messages come in reverse chronological order, reverse them
	var sb strings.Builder
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		author := "Unknown"
		if msg.Author != nil {
			author = msg.Author.Username
		}
		ts := msg.Timestamp
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", ts, author, msg.Content))

		// Show attachments
		for _, att := range msg.Attachments {
			sb.WriteString(fmt.Sprintf("  📎 %s (%s)\n", att.Filename, att.URL))
		}

		// Show embeds summary
		for _, embed := range msg.Embeds {
			if embed.Title != "" {
				sb.WriteString(fmt.Sprintf("  📦 Embed: %s\n", embed.Title))
			}
		}
	}
	return sb.String()
}
