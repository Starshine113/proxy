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

import "errors"

// ErrorNoSystem is returned when the user has no system
var ErrorNoSystem = errors.New("user doesn't have a system")

// CheckSystem checks if the user has a system
func (ctx *Ctx) CheckSystem() (err error) {
	if !ctx.Database.HasSystem(ctx.Author.ID) {
		ctx.Sendf("%v You do not have a system registered with %v. To create one, use `%vsystem new`.", ErrorEmoji, ctx.BotUser.Username, ctx.Bot.Prefix)
		return ErrorNoSystem
	}

	return nil
}
