package member

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
	"github.com/Starshine113/proxy/etc"
	"strings"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func avatar(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if err = ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	members, err := ctx.Database.GetAccountMembers(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	member := ctx.Args[0]

	var m *db.Member

	for _, mem := range members {
		if member == mem.ID.String() {
			m = mem
			break
		}
		if strings.ToLower(member) == strings.ToLower(mem.Name) {
			m = mem
			break
		}
	}

	if m == nil {
		_, err = ctx.Sendf("%v Member `%v` not found.\n**Note:** if the member has spaces in their name, you need to use the ID.", router.ErrorEmoji, member)
		return
	}

	if len(ctx.Args) == 1 {
		if len(ctx.Message.Attachments) == 0 {
			return viewAvatar(m, ctx)
		}
		return changeAvatar(m, ctx)
	}

	return changeAvatar(m, ctx)
}

func changeAvatar(m *db.Member, ctx *router.Ctx) (err error) {
	if len(ctx.Args) > 1 {
		member, err := ctx.ParseMember(strings.Join(ctx.Args[1:], " "))
		if err == nil {
			err = ctx.Database.SetAvatar(m.ID.String(), member.User.AvatarURL("256"))
			if err != nil {
				return ctx.CommandError(err)
			}
			_, err = ctx.Sendf("%v Member avatar changed to %v's avatar.\nNote that if %v changes their avatar, this member's avatar will need to be reset.", router.SuccessEmoji, member.Mention(), member.User.Username)
			return err
		}

		// it's not a member (or the fetch failed), so try if it's an image
		if etc.HasAnySuffix(ctx.Args[1], ".jpg", ".jpeg", ".png", ".gif", ".webp") {
			err = ctx.Database.SetAvatar(m.ID.String(), ctx.Args[1])
			if err != nil {
				return ctx.CommandError(err)
			}
			_, err = ctx.Send(&discordgo.MessageSend{
				Content: fmt.Sprintf("%v Member avatar changed to the image at the given URL.", router.SuccessEmoji),
				Embed: &discordgo.MessageEmbed{
					Title: fmt.Sprintf("%v's avatar", m.Name),
					Image: &discordgo.MessageEmbedImage{
						URL: ctx.Args[1],
					},
				},
			})
			return err
		}
	}

	if len(ctx.Message.Attachments) == 0 {
		return nil
	}

	a := ctx.Message.Attachments[0]

	if etc.HasAnySuffix(a.URL, ".jpg", ".jpeg", ".png", ".gif", ".webp") {
		err = ctx.Database.SetAvatar(m.ID.String(), a.URL)
		if err != nil {
			return ctx.CommandError(err)
		}
		_, err = ctx.Send(&discordgo.MessageSend{
			Content: fmt.Sprintf("%v Member avatar changed to the attached image.\nNote that if this message is deleted, the member's avatar will need to be reset.", router.SuccessEmoji),
			Embed: &discordgo.MessageEmbed{
				Title: fmt.Sprintf("%v's avatar", m.Name),
				Image: &discordgo.MessageEmbedImage{
					URL: a.URL,
				},
			},
		})
		return err
	}

	return
}

func viewAvatar(m *db.Member, ctx *router.Ctx) (err error) {
	if m.AvatarURL == "" {
		_, err = ctx.Sendf("%v has no avatar set.", m.Name)
		return err
	}

	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%v's avatar", m.Name),
		Description: fmt.Sprintf("To clear, use `%vmember avatar %s clear`.", ctx.Bot.Prefix, m.ID),
		Image: &discordgo.MessageEmbedImage{
			URL: m.AvatarURL,
		},
	})
	return err
}
