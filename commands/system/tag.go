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
)

func tag(ctx *router.Ctx) (err error) {
	tag := strings.Join(ctx.Args, " ")

	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	if tag == "" {
		if s.Tag == "" {
			_, err = ctx.Sendf("%v You do not currently have a system tag set. To set one, use `%vsystem tag <new tag>`.", router.ErrorEmoji, ctx.Bot.Prefix)
			return
		}
		_, err = ctx.Sendf("Your current system tag is `%v`. To change it, use `%vsystem tag <new tag>`.", s.Tag, ctx.Bot.Prefix)
		return
	}

	if len(tag) > 100 {
		_, err = ctx.Sendf("%v Your new system tag would be too long (%v > 100 characters).", router.ErrorEmoji, len(tag))
		return
	}

	err = ctx.Database.SetTag(s.ID.String(), tag)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Succesfully changed your system tag to `%v`.", router.SuccessEmoji, tag)
	return
}
