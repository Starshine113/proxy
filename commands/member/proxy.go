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

func proxy(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if ctx.CheckMinArgs(2); err != nil {
		return ctx.CommandError(err)
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	members, err := ctx.Database.GetSystemMembers(s.ID.String())
	if err != nil {
		return ctx.CommandError(err)
	}

	member := ctx.Args[0]
	proxy := strings.Join(ctx.Args[1:], " ")

	if proxy == "text" {
		_, err = ctx.Sendf("%v Proxy can't be empty.", router.ErrorEmoji)
		return
	} else if !strings.Contains(proxy, "text") {
		_, err = ctx.Sendf("%v Invalid proxy tag supplied.\nTo set a proxy, pretend to proxy the message \"text\"; for example, inputting `[text]` will proxy any message surrounded by square brackets as that member.", router.ErrorEmoji)
		return
	}

	tags := strings.Split(proxy, "text")

	newProxy := db.ProxyTag{}

	if len(tags) > 1 {
		newProxy.Prefix = tags[0]
		newProxy.Suffix = tags[1]
	} else {
		if strings.HasPrefix(proxy, "text") {
			newProxy.Suffix = tags[0]
		} else {
			newProxy.Prefix = tags[0]
		}
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

	if m == nil {
		_, err = ctx.Sendf("%v Member `%v` not found.\n**Note:** if the member has spaces in their name, you need to use the ID.", router.ErrorEmoji, member)
		return
	}

	err = ctx.Database.SetProxy(m.ID.String(), newProxy.Prefix, newProxy.Suffix)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Member proxy updated. This member will now be proxied with the tag `%vtext%v.`", router.SuccessEmoji, newProxy.Prefix, newProxy.Suffix)
	return
}
