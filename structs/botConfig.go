package structs

// BotConfig holds the bot's configuration
type BotConfig struct {
	Auth struct {
		Token       string `toml:"token"`
		DatabaseURL string `toml:"database_url"`
	} `toml:"auth"`
	Bot struct {
		Prefixes   []string `toml:"prefixes"`
		BotOwners  []string `toml:"bot_owners"`
		LogWebhook string   `toml:"log_webhook"`
	} `toml:"bot"`
}
