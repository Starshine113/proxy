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

// AddReactionHandlerOnce adds a reaction handler function that is only called once
func (ctx *Ctx) AddReactionHandlerOnce(messageID, reaction string, f func(ctx *Ctx)) func() {
	returnFunc := ctx.Bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID != ctx.Message.Author.ID || r.ChannelID != ctx.Channel.ID || r.MessageID != messageID || r.MessageReaction.Emoji.APIName() != reaction {
			return
		}
		f(ctx)

		if v, e := ctx.Bot.Handlers.Get(messageID + reaction); e == nil {
			v.(func())()
			ctx.Bot.Handlers.Remove(messageID + reaction)
		}

		return
	})
	ctx.Bot.Handlers.Set(messageID+reaction, returnFunc)
	return returnFunc
}

// AddReactionHandler adds a reaction handler function
func (ctx *Ctx) AddReactionHandler(messageID, reaction string, f func(ctx *Ctx)) func() {
	returnFunc := ctx.Bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID != ctx.Message.Author.ID || r.ChannelID != ctx.Channel.ID || r.MessageID != messageID || r.MessageReaction.Emoji.APIName() != reaction {
			return
		}
		f(ctx)

		return
	})
	ctx.Bot.Handlers.Set(messageID+reaction, returnFunc)
	return returnFunc
}

// AddYesNoHandler reacts with ✅ and ❌, and runs one of two functions depending on which one is used
func (ctx *Ctx) AddYesNoHandler(messageID string, yesFunc, noFunc func(ctx *Ctx)) func() {
	ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, messageID, "✅")
	ctx.Bot.Session.MessageReactionAdd(ctx.Channel.ID, messageID, "❌")

	returnFunc := ctx.Bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID != ctx.Message.Author.ID || r.ChannelID != ctx.Channel.ID || r.MessageID != messageID {
			return
		}

		switch r.MessageReaction.Emoji.APIName() {
		case "✅":
			yesFunc(ctx)
		case "❌":
			noFunc(ctx)
		default:
			return
		}

		if v, e := ctx.Bot.Handlers.Get(messageID); e == nil {
			v.(func())()
			ctx.Bot.Handlers.Remove(messageID)
		}

		return
	})
	ctx.Bot.Handlers.Set(messageID, returnFunc)
	return returnFunc
}
