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
	"strings"

	"github.com/Starshine113/proxy/router"
)

func new(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	name := strings.Join(ctx.Args, " ")
	if name == "" {
		_, err = ctx.Sendf("%v You must pass a name.", router.ErrorEmoji)
		return
	}

	if len(name) > 100 {
		_, err = ctx.Sendf("%v The name you gave is too long (%v > 100 characters).", router.ErrorEmoji, len(name))
		return
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	m, err := ctx.Database.NewMember(s.ID.String(), name)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Member \"%v\" (`%v`) registered!\nTo get started proxying this member, use the command `%vmember proxy %v`.", router.SuccessEmoji, m.Name, m.ID, ctx.Bot.Prefix, m.ID)
	return
}
