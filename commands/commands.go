package commands

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
	"github.com/Starshine113/proxy/commands/member"
	"github.com/Starshine113/proxy/commands/static"
	"github.com/Starshine113/proxy/commands/system"
	"github.com/Starshine113/proxy/router"
)

// Init ...
func Init(r *router.Router) {
	static.Init(r)
	system.Init(r)
	member.Init(r)
}
