package bot

// Discord proxy bot
// Copyright (C) 2020  Starshine System

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
