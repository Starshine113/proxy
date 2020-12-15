package router

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

import "github.com/bwmarrin/discordgo"

// AddMessageHandler adds a listener for a message from a specific user
func (ctx *Ctx) AddMessageHandler(messageID string, f func(ctx *Ctx, m *discordgo.MessageCreate)) func() {
	returnFunc := ctx.Bot.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID != ctx.Message.Author.ID || m.ChannelID != ctx.Channel.ID {
			return
		}
		f(ctx, m)

		if v, e := ctx.Bot.Handlers.Get(messageID); e == nil {
			v.(func())()
			ctx.Bot.Handlers.Remove(messageID)
		}

		return
	})
	ctx.Bot.Handlers.Set(messageID, returnFunc)
	return returnFunc
}
