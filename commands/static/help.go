package static

import (
	"fmt"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func help(ctx *router.Ctx) (err error) {
	b := router.NewEmbedBuilder(fmt.Sprintf("%v help", ctx.BotUser.Username), ctx.BotUser.Username, ctx.BotUser.AvatarURL("128"), 0x21a1a8)

	b.Add("Getting started", fmt.Sprintf("To get started using %v, use the following commands, use `%vsystem new`. From there, follow the commands below.", ctx.BotUser.Username, ctx.Bot.Prefix), []*discordgo.MessageEmbedField{{
		Name: "Creating a member",
		Value: fmt.Sprintf("**1.** `%vmember new John` - Add a new member to your system\n**2.** `%vmember proxy John [text]` - Set up [square brackets] as proxy tags\n**3.** You're done! You can now type [a message in brackets] and it'll be proxied appropriately.\n**5.** Optionally, you may set an avatar from the URL of an image with `%vmember avatar John [link to image]`, or from a file by typing `%vmember avatar John` and sending the message with an attached image.", ctx.Bot.Prefix, ctx.Bot.Prefix, ctx.Bot.Prefix, ctx.Bot.Prefix),
	}, {
		Name: "Setting up your system",
		Value: fmt.Sprintf("**1.** `%vs rename New Name` - change your system name\n**2.** `%vs tag | New Tag` - change your system tag, which will show up after a proxied member's name", ctx.Bot.Prefix, ctx.Bot.Prefix),
	}, {
		Name: "Invite the bot",
		Value: fmt.Sprintf("To invite the bot to your server, use [this](%v) link.", ctx.Invite()),
	}})

	_, err = ctx.PagedEmbed(b.Build())
	return err
}