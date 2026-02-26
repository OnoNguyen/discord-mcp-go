package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/houseofdoge/discord-mcp-go/config"
)

type Client struct {
	Session  *discordgo.Session
	GuildID  string // default guild ID from config
}

func NewClient(cfg *config.Config) (*Client, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("creating discord session: %w", err)
	}

	// Validate the token by fetching the bot user
	_, err = session.User("@me")
	if err != nil {
		return nil, fmt.Errorf("invalid discord token: %w", err)
	}

	return &Client{
		Session: session,
		GuildID: cfg.GuildID,
	}, nil
}

// ResolveGuildID returns the provided guildID if non-empty, otherwise falls back to the default.
func (c *Client) ResolveGuildID(guildID string) (string, error) {
	if guildID != "" {
		return guildID, nil
	}
	if c.GuildID != "" {
		return c.GuildID, nil
	}
	return "", fmt.Errorf("guildId is required (no default DISCORD_GUILD_ID configured)")
}
