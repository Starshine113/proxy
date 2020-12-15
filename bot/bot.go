package bot

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/structs"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Bot ...
type Bot struct {
	Db       *db.Db
	Session  *discordgo.Session
	Config   *structs.BotConfig
	Sugar    *zap.SugaredLogger
	Handlers *ttlcache.Cache

	Prefix string
}

// New returns a new bot instance
func New(db *db.Db, s *discordgo.Session, c *structs.BotConfig, l *zap.SugaredLogger) *Bot {
	handlerMap := ttlcache.NewCache()
	handlerMap.SetCacheSizeLimit(10000)
	handlerMap.SetTTL(15 * time.Minute)
	handlerMap.SetExpirationCallback(func(key string, value interface{}) {
		value.(func())()
	})

	b := &Bot{Db: db, Session: s, Config: c, Sugar: l, Handlers: handlerMap, Prefix: c.Bot.Prefixes[0]}

	return b
}
