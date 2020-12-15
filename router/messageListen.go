package router

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
