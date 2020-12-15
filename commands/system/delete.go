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
	"strings"

	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func delete(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Are you sure you want to delete your system? If so, reply with your system ID (`%s`).\n*This action is irreversible!*", router.WarnEmoji, s.ID)

	ctx.AddMessageHandler(ctx.Message.ID, func(ctx *router.Ctx, m *discordgo.MessageCreate) {
		c := strings.TrimSpace(m.Content)

		if c != s.ID.String() {
			_, err = ctx.Sendf("%v System deletion cancelled. Note that you must reply with your system ID (`%s`) *verbatim*.", router.ErrorEmoji, s.ID)
			return
		}
		err = ctx.Database.DeleteSystem(c)
		if err != nil {
			ctx.CommandError(err)
			return
		}
		_, err = ctx.Sendf("%v System deleted.", router.SuccessEmoji)
		return
	})

	return
}
