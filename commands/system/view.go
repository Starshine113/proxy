package system

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
	"strings"
	"time"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func view(ctx *router.Ctx) (err error) {
	var s *db.System

	if len(ctx.Args) > 0 {
		if !ctx.Database.HasSystem(ctx.Author.ID) {
			_, err = ctx.Sendf("%v You do not have a system registered with %v. To create one, use `%vsystem new`.", router.ErrorEmoji, ctx.BotUser.Username, ctx.Bot.Prefix)
			return
		}

		s, err = ctx.Database.GetUserSystem(ctx.Author.ID)
		if err != nil {
			return ctx.CommandError(err)
		}
	}

	u, err := ctx.Database.GetSystemUsers(s.ID.String())
	if err != nil {
		return ctx.CommandError(err)
	}

	users := strings.Join(router.PrintfAll("<@%v>", u), "\n")

	fields := make([]*discordgo.MessageEmbedField, 0)

	if s.Tag != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Tag",
			Value:  s.Tag,
			Inline: false,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Linked accounts",
		Value:  users,
		Inline: true,
	})

	members, err := ctx.Database.GetSystemMembers(s.ID.String())

	var m string
	if len(members) > 0 {
		m = fmt.Sprintf("(see `%vsystem list %v`)", ctx.Bot.Prefix, s.ID)
	} else {
		var ownSystem bool
		for _, id := range u {
			if id == ctx.Author.ID {
				ownSystem = true
			}
		}
		if ownSystem {
			m = fmt.Sprintf("Add one with `%vmember new`!", ctx.Bot.Prefix)
		} else {
			m = "No members."
		}
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("Members (%v)", len(members)),
		Value:  m,
		Inline: true,
	}, &discordgo.MessageEmbedField{
		Name:   "System ID",
		Value:  "```" + s.ID.String() + "```",
		Inline: false,
	})

	embed := &discordgo.MessageEmbed{
		Title:  s.Name,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Created on",
		},
		Timestamp: s.Created.UTC().Format(time.RFC3339),
	}

	_, err = ctx.Send(embed)
	return
}
