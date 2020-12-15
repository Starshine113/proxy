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

import (
	"github.com/bwmarrin/discordgo"
)

// PagedEmbed ...
func (ctx *Ctx) PagedEmbed(embeds []*discordgo.MessageEmbed) (msg *discordgo.Message, err error) {
	if len(embeds) == 1 {
		return ctx.Send(embeds[0])
	}
	msg, err = ctx.Send(embeds[0])
	if err != nil {
		return
	}
	if err = ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, "❌"); err != nil {
		return
	}
	if err = ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, "⏪"); err != nil {
		return
	}
	if err = ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, "⬅️"); err != nil {
		return
	}
	if err = ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, "➡️"); err != nil {
		return
	}
	if err = ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, "⏩"); err != nil {
		return
	}

	ctx.AdditionalParams["page"] = 0
	ctx.AdditionalParams["embeds"] = embeds

	ctx.AddReactionHandler(msg.ID, "⬅️", func(ctx *Ctx) {
		page := ctx.AdditionalParams["page"].(int)
		embeds := ctx.AdditionalParams["embeds"].([]*discordgo.MessageEmbed)

		if ctx.Message.GuildID != "" {
			ctx.Bot.Session.MessageReactionRemove(ctx.Channel.ID, msg.ID, "⬅️", ctx.Author.ID)
		}

		if page == 0 {
			return
		}
		ctx.Edit(msg, embeds[page-1])
		ctx.AdditionalParams["page"] = page - 1
	})

	ctx.AddReactionHandler(msg.ID, "➡️", func(ctx *Ctx) {
		page := ctx.AdditionalParams["page"].(int)
		embeds := ctx.AdditionalParams["embeds"].([]*discordgo.MessageEmbed)

		if ctx.Message.GuildID != "" {
			ctx.Bot.Session.MessageReactionRemove(ctx.Channel.ID, msg.ID, "➡️", ctx.Author.ID)
		}

		if page >= len(embeds)-1 {
			return
		}
		ctx.Edit(msg, embeds[page+1])
		ctx.AdditionalParams["page"] = page + 1
	})

	ctx.AddReactionHandler(msg.ID, "⏪", func(ctx *Ctx) {
		embeds := ctx.AdditionalParams["embeds"].([]*discordgo.MessageEmbed)

		if ctx.Message.GuildID != "" {
			ctx.Bot.Session.MessageReactionRemove(ctx.Channel.ID, msg.ID, "⏪", ctx.Author.ID)
		}

		ctx.Edit(msg, embeds[0])
		ctx.AdditionalParams["page"] = 0
	})

	ctx.AddReactionHandler(msg.ID, "⏩", func(ctx *Ctx) {
		embeds := ctx.AdditionalParams["embeds"].([]*discordgo.MessageEmbed)

		if ctx.Message.GuildID != "" {
			ctx.Bot.Session.MessageReactionRemove(ctx.Channel.ID, msg.ID, "⏩", ctx.Author.ID)
		}

		ctx.Edit(msg, embeds[len(embeds)-1])
		ctx.AdditionalParams["page"] = len(embeds) - 1
	})

	return msg, err
}
