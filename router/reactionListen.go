package router

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
