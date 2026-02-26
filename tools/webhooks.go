package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListWebhooksInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID to list webhooks from"`
}

type CreateWebhookInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID to create the webhook in"`
	Name      string `json:"name" jsonschema:"Name for the webhook"`
}

type DeleteWebhookInput struct {
	WebhookID string `json:"webhookId" jsonschema:"ID of the webhook to delete"`
}

type SendWebhookMessageInput struct {
	WebhookID    string `json:"webhookId" jsonschema:"Webhook ID"`
	WebhookToken string `json:"webhookToken" jsonschema:"Webhook token"`
	Content      string `json:"content" jsonschema:"Message content to send"`
}

func RegisterWebhookTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_webhooks",
		Description: "List webhooks in a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListWebhooksInput) (*mcp.CallToolResult, any, error) {
		webhooks, err := dc.Session.ChannelWebhooks(input.ChannelID)
		if err != nil {
			return errResult(fmt.Errorf("fetching webhooks: %w", err)), nil, nil
		}

		if len(webhooks) == 0 {
			return textResult("No webhooks found."), nil, nil
		}

		var sb strings.Builder
		for _, wh := range webhooks {
			sb.WriteString(fmt.Sprintf("🔗 %s (ID: %s, Token: %s)\n", wh.Name, wh.ID, wh.Token))
		}
		return textResult(sb.String()), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_webhook",
		Description: "Create a new webhook in a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateWebhookInput) (*mcp.CallToolResult, any, error) {
		wh, err := dc.Session.WebhookCreate(input.ChannelID, input.Name, "")
		if err != nil {
			return errResult(fmt.Errorf("creating webhook: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Webhook '%s' created (ID: %s, Token: %s)", wh.Name, wh.ID, wh.Token)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_webhook",
		Description: "Delete a Discord webhook",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteWebhookInput) (*mcp.CallToolResult, any, error) {
		err := dc.Session.WebhookDelete(input.WebhookID)
		if err != nil {
			return errResult(fmt.Errorf("deleting webhook: %w", err)), nil, nil
		}
		return textResult(fmt.Sprintf("Webhook %s deleted", input.WebhookID)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_webhook_message",
		Description: "Send a message using a Discord webhook",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SendWebhookMessageInput) (*mcp.CallToolResult, any, error) {
		_, err := dc.Session.WebhookExecute(input.WebhookID, input.WebhookToken, false, &discordgo.WebhookParams{
			Content: input.Content,
		})
		if err != nil {
			return errResult(fmt.Errorf("sending webhook message: %w", err)), nil, nil
		}
		return textResult("Webhook message sent successfully"), nil, nil
	})
}
