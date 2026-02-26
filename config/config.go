package config

import (
	"fmt"
	"os"
)

type Config struct {
	DiscordToken string
	GuildID      string // optional default guild
}

func Load() (*Config, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable is required")
	}

	return &Config{
		DiscordToken: token,
		GuildID:      os.Getenv("DISCORD_GUILD_ID"),
	}, nil
}
