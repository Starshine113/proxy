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

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
)

func displayName(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if ctx.CheckMinArgs(2); err != nil {
		return ctx.CommandError(err)
	}

	members, err := ctx.Database.GetAccountMembers(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	member := ctx.Args[0]
	name := strings.Join(ctx.Args[1:], " ")

	if len(name) == 0 {
		_, err = ctx.Sendf("%v You must specify a new display name.", router.ErrorEmoji)
		return err
	} else if len(name) > 80 {
		_, err = ctx.Sendf("%v The display name you gave was too long (%v > 80 characters).", router.ErrorEmoji, len(name))
		return err
	}

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

	err = ctx.Database.SetDisplayName(m.ID.String(), name)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Member display name updated. This member will now be proxied with the name `%v.`", router.SuccessEmoji, name)
	return
}
