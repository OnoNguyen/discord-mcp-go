package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/houseofdoge/discord-mcp-go/config"
	"github.com/houseofdoge/discord-mcp-go/discord"
	"github.com/houseofdoge/discord-mcp-go/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dc, err := discord.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Discord client: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "discord-mcp-go",
		Version: "v1.0.0",
	}, nil)

	// Register all tool groups
	tools.RegisterServerInfoTools(server, dc)
	tools.RegisterMessageTools(server, dc)
	tools.RegisterReactionTools(server, dc)
	tools.RegisterChannelTools(server, dc)
	tools.RegisterThreadTools(server, dc)
	tools.RegisterCategoryTools(server, dc)
	tools.RegisterUserTools(server, dc)
	tools.RegisterRoleTools(server, dc)
	tools.RegisterWebhookTools(server, dc)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
