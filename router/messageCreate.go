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

// MessageCreate ...
func (r *Router) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// if message was sent by a bot return; not only to ignore bots, but also to make sure PluralKit users don't trigger commands twice.
	if m.Author.Bot {
		return
	}

	ctx, err := Context(r.Bot.Config.Bot.Prefixes, m.Content, m, r.Bot)
	if err != nil {
		r.Bot.Sugar.Errorf("Error getting context for %v: %v", m.ID, err)
		return
	}

	// check if the message might be a command
	if ctx.MatchPrefix() {
		r.execute(ctx)
		return
	}

	r.Proxy.MessageCreate(s, m)
}

func (r *Router) execute(ctx *Ctx) {
	err := r.Execute(ctx)
	if err != nil {
		r.Bot.Sugar.Errorf("Error running command %v: %v", ctx.Message.ID, err)
	}
}
