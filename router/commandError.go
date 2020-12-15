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
	"fmt"
	"time"

	"codeberg.org/eviedelta/dwhook"
	"github.com/bwmarrin/discordgo"
)

// CommandError sends an error message and optionally returns an error for logging purposes
func (ctx *Ctx) CommandError(err error) error {
	switch err {
	case ErrorNoDMs, ErrorMissingBotOwner, ErrorMissingManagerPerms:
		ctx.React(WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " You are not allowed to use this command:\n> " + err.Error(),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return msgErr
	case ErrorNotEnoughArgs:
		ctx.React(WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " Command call was missing arguments.",
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return msgErr
	case ErrorTooManyArgs:
		ctx.React(WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " Command call has too many arguments.",
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return msgErr
	}
	switch err.(type) {
	case *discordgo.RESTError:
		e := err.(*discordgo.RESTError)
		if e.Message != nil {
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "REST error occurred",
				Description: fmt.Sprintf("```%v ```", e.Message.Message),
				Fields: []*discordgo.MessageEmbedField{{
					Name:   "Raw",
					Value:  string(e.ResponseBody),
					Inline: false,
				}},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Error code: %v", e.Message.Code),
				},
				Color:     0xbf1122,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
		} else {
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "REST error occurred",
				Description: fmt.Sprintf("```%v```", e.ResponseBody),
				Color:       0xbf1122,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			})
		}
		return err
	default:
		ctx.Bot.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, ErrorEmoji)

		embed := &discordgo.MessageEmbed{
			Title:       "Internal error occured",
			Description: fmt.Sprintf("```%v```\nIf this error persists, please contact the bot developer.", err.Error()),
			Color:       0xbf1122,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}

		config := ctx.Bot.Config

		if config.Bot.LogWebhook != "" {
			msg := dwhook.Message{
				Content:   fmt.Sprintf("> An internal error occured in %v (%v) of guild %v\n> Triggered by %v (%v/%v):", ctx.Channel.ID, ctx.Channel.Mention(), ctx.Channel.GuildID, ctx.Author.String(), ctx.Author.Mention(), ctx.Author.ID),
				Username:  ctx.BotUser.Username + " Error",
				AvatarURL: ctx.BotUser.AvatarURL("256"),
				Embeds: []dwhook.Embed{{
					Color:       0xbf1122,
					Description: fmt.Sprintf("```%v```", err.Error()),
					Fields: []dwhook.EmbedField{{
						Name:  "Command",
						Value: fmt.Sprintf("**Command**: `%v`\n**Arguments**: `%v`", ctx.Command, ctx.Args),
					}},
					Footer: dwhook.EmbedFooter{
						Text: "Triggering message ID: " + ctx.Author.ID,
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				}},
			}
			dwhook.SendTo(config.Bot.LogWebhook, msg)
		}

		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Embed: embed,
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return msgErr
	}
}
