package router

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Send a message to the context channel
func (ctx *Ctx) Send(arg interface{}) (message *discordgo.Message, err error) {
	message, err = ctx.SendNoAddXHandler(arg)
	if err != nil {
		return
	}
	ctx.AddReactionHandlerOnce(message.ID, "❌", func(ctx *Ctx) {
		err := ctx.Bot.Session.ChannelMessageDelete(ctx.Channel.ID, message.ID)
		if err != nil {
			ctx.Bot.Sugar.Errorf("Error deleting message %v: %v", message.ID, err)
		}
	})
	return
}

// SendNoAddXHandler sends a message without adding a handler for :x:
func (ctx *Ctx) SendNoAddXHandler(arg interface{}) (message *discordgo.Message, err error) {
	switch arg.(type) {
	case string:
		message, err = ctx.Bot.Session.ChannelMessageSend(ctx.Message.ChannelID, arg.(string))
	case *discordgo.MessageEmbed:
		message, err = ctx.Bot.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, arg.(*discordgo.MessageEmbed))
	case *discordgo.MessageSend:
		message, err = ctx.Bot.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, arg.(*discordgo.MessageSend))
	default:
		err = errors.New("don't know what to do with that type")
	}
	return message, err
}

// Sendf sends a fmt.Sprintf-like input string
func (ctx *Ctx) Sendf(format string, args ...interface{}) (msg *discordgo.Message, err error) {
	msg, err = ctx.SendfNoAddXHandler(format, args...)
	ctx.AddReactionHandlerOnce(msg.ID, "❌", func(ctx *Ctx) {
		ctx.Bot.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
	})
	return
}

// SendfNoAddXHandler ...
func (ctx *Ctx) SendfNoAddXHandler(format string, args ...interface{}) (msg *discordgo.Message, err error) {
	return ctx.Bot.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(format, args...))
}

// Editf edits a message with Sendf-like syntax
func (ctx *Ctx) Editf(message *discordgo.Message, format string, args ...interface{}) (msg *discordgo.Message, err error) {
	return ctx.Bot.Session.ChannelMessageEdit(message.ChannelID, message.ID, fmt.Sprintf(format, args...))
}

// Embed sends the input as an embed
func (ctx *Ctx) Embed(title, description string, color int) (msg *discordgo.Message, err error) {
	msg, err = ctx.EmbedNoXHandler(title, description, color)
	ctx.AddReactionHandlerOnce(msg.ID, "❌", func(ctx *Ctx) {
		ctx.Bot.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
	})
	return
}

// EmbedNoXHandler ...
func (ctx *Ctx) EmbedNoXHandler(title, description string, color int) (msg *discordgo.Message, err error) {
	if color == 0 {
		color = 0x21a1a8
	}
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	return ctx.Bot.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
}

// Embedf sends a fmt.Sprintf-like input string, in an embed
func (ctx *Ctx) Embedf(title, format string, args ...interface{}) (msg *discordgo.Message, err error) {
	msg, err = ctx.EmbedfNoXHandler(title, format, args...)
	ctx.AddReactionHandlerOnce(msg.ID, "❌", func(ctx *Ctx) {
		ctx.Bot.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
	})
	return
}

// EmbedfNoXHandler ...
func (ctx *Ctx) EmbedfNoXHandler(title, format string, args ...interface{}) (msg *discordgo.Message, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf(format, args...),
		Color:       0x21a1a8,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	return ctx.Bot.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
}

// EditEmbedf edits an embed with Embedf syntax
func (ctx *Ctx) EditEmbedf(message *discordgo.Message, title, format string, args ...interface{}) (msg *discordgo.Message, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf(format, args...),
		Color:       0x21a1a8,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	return ctx.Bot.Session.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)
}

// Edit a message
func (ctx *Ctx) Edit(message *discordgo.Message, arg interface{}) (msg *discordgo.Message, err error) {
	switch arg.(type) {
	case string:
		msg, err = ctx.Bot.Session.ChannelMessageEdit(message.ChannelID, message.ID, arg.(string))
	case *discordgo.MessageEmbed:
		msg, err = ctx.Bot.Session.ChannelMessageEditEmbed(message.ChannelID, message.ID, arg.(*discordgo.MessageEmbed))
	case *discordgo.MessageEdit:
		edit := arg.(*discordgo.MessageEdit)
		edit.ID = message.ID
		edit.Channel = message.ChannelID
		msg, err = ctx.Bot.Session.ChannelMessageEditComplex(edit)
	default:
		err = errors.New("don't know what to do with that type")
	}
	return msg, err
}

// React reacts to the triggering message
func (ctx *Ctx) React(emoji string) (err error) {
	return ctx.Bot.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, emoji)
}

// TriggerTyping triggers typing in the channel the command was invoked in.
func (ctx *Ctx) TriggerTyping() (err error) {
	return ctx.Bot.Session.ChannelTyping(ctx.Message.ChannelID)
}
