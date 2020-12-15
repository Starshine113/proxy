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
	"strings"
	"time"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func view(ctx *router.Ctx) (err error) {
	member := strings.Join(ctx.Args, " ")
	var m *db.Member
	if len(member) == 0 {
		_, err = ctx.Sendf("%v You must specify a member (name or ID).", router.ErrorEmoji)
		return
	}
	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err == nil {
		members, err := ctx.Database.GetSystemMembers(s.ID.String())
		if err != nil {
			return ctx.CommandError(err)
		}
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
	}

	if m == nil {
		_, err = ctx.Sendf("%v Member `%v` not found.\n**Note:** if the member has spaces in their name, you need to use the ID.", router.ErrorEmoji, member)
		return
	}

	m, err = ctx.Database.MemberWithMsgCount(m.ID.String())
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Send(memberCard(s, m))
	return err
}

func memberCard(s *db.System, m *db.Member) *discordgo.MessageEmbed {
	title := m.Name
	if s.Name != "" {
		title += " (" + s.Name + ")"
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	if m.DisplayName != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Display name",
			Value:  m.DisplayName,
			Inline: true,
		})
	}

	if m.Prefix != "" || m.Suffix != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Proxy",
			Value:  fmt.Sprintf("`%vtext%v`", m.Prefix, m.Suffix),
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Message count",
		Value:  fmt.Sprint(m.MessageCount),
		Inline: true,
	}, &discordgo.MessageEmbedField{
		Name:   "** **",
		Value:  fmt.Sprintf("Member ID: `%s`\nSystem ID: `%s`", m.ID, s.ID),
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title: title,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: m.AvatarURL,
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Created at",
		},
		Timestamp: m.Created.UTC().Format(time.RFC3339),
	}
}
