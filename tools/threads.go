package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListThreadsInput struct {
	ChannelID string `json:"channelId" jsonschema:"Channel ID to list threads from"`
}

type ReadThreadMessagesInput struct {
	ThreadID string `json:"threadId" jsonschema:"Thread ID to read messages from"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Number of messages to fetch (default 100, max 100)"`
}

func RegisterThreadTools(server *mcp.Server, dc *discord.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_threads",
		Description: "List active threads in a Discord channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListThreadsInput) (*mcp.CallToolResult, any, error) {
		threads, err := dc.Session.ThreadsActive(input.ChannelID)
		if err != nil {
			return errResult(fmt.Errorf("fetching threads: %w", err)), nil, nil
		}

		if len(threads.Threads) == 0 {
			return textResult("No active threads found."), nil, nil
		}

		var sb strings.Builder
		for _, t := range threads.Threads {
			sb.WriteString(fmt.Sprintf("🧵 %s (ID: %s, Messages: %d)\n", t.Name, t.ID, t.MessageCount))
		}
		return textResult(sb.String()), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_thread_messages",
		Description: "Read messages from a Discord thread. Returns messages in chronological order.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ReadThreadMessagesInput) (*mcp.CallToolResult, any, error) {
		limit := input.Limit
		if limit <= 0 || limit > 100 {
			limit = 100
		}

		messages, err := dc.Session.ChannelMessages(input.ThreadID, limit, "", "", "")
		if err != nil {
			return errResult(fmt.Errorf("fetching thread messages: %w", err)), nil, nil
		}

		return textResult(formatMessages(messages)), nil, nil
	})
}
