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

import (
	"github.com/Starshine113/proxy/etc"
	"regexp"
	"strings"
)

// MatchPrefix checks if the message matched any prefix
func (ctx *Ctx) MatchPrefix() bool {
	return etc.HasAnyPrefix(strings.ToLower(ctx.Message.Content), ctx.Bot.Config.Bot.Prefixes...)
}

// Match checks if any of the given command aliases match
func (ctx *Ctx) Match(cmds ...string) bool {
	for _, cmd := range cmds {
		if strings.ToLower(ctx.Command) == strings.ToLower(cmd) {
			return true
		}
	}
	return false
}

// MatchRegexp checks if the command matches the given regex
func (ctx *Ctx) MatchRegexp(re *regexp.Regexp) bool {
	if re == nil {
		return false
	}
	return re.MatchString(strings.ToLower(ctx.Command))
}
