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
	"sort"
	"strings"
	"time"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

type mList []*db.Member

func (c mList) Len() int      { return len(c) }
func (c mList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c mList) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{c[i].Name, c[j].Name})
}

func list(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	members := make(mList, 0)

	members, err = ctx.Database.GetSystemMembers(s.ID.String())
	if err != nil {
		return ctx.CommandError(err)
	}

	sort.Sort(members)

	memberSlices := make([][]*db.Member, 0)

	for i := 0; i < len(members); i += 10 {
		end := i + 10

		if end > len(members) {
			end = len(members)
		}

		memberSlices = append(memberSlices, members[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)
	title := fmt.Sprintf("Members of `%s`", s.ID)
	if s.Name != "" {
		title = fmt.Sprintf("Members of %v (`%s`)", s.Name, s.ID)
	}

	for i, c := range memberSlices {
		x := make([]string, 0)
		for _, member := range c {
			x = append(x, fmt.Sprintf("[`%s`] %v", member.ID, member.Name))
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title: title,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(memberSlices)),
			},
			Timestamp:   time.Now().Format(time.RFC3339),
			Description: strings.Join(x, "\n"),
		})
	}

	_, err = ctx.PagedEmbed(embeds)
	return
}
