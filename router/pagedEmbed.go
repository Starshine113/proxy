package router

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
